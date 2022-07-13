package iaas

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/DataWorkbench/common/gtrace"
	"github.com/DataWorkbench/common/web/ghttp"
	"github.com/DataWorkbench/glog"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"
)

var traceComponentTag = opentracing.Tag{Key: string(ext.Component), Value: "IaasClient"}

const (
	iso8601            = "2006-01-02T15:04:05Z"
	signatureAlgorithm = "HmacSHA256"
	signatureVersion   = "1"
)

// Config represents the iaas api config.
type Config struct {
	Zone            string `json:"zone"              yaml:"zone"              env:"ZONE"                validate:"required"`
	Host            string `json:"host"              yaml:"host"              env:"HOST"                validate:"required"`
	Port            int    `json:"port"              yaml:"port"              env:"PORT"                validate:"required"`
	Protocol        string `json:"protocol"          yaml:"protocol"          env:"PROTOCOL"            validate:"required"`
	Timeout         int    `json:"timeout"           yaml:"timeout"           env:"TIMEOUT,default=600" validate:"required"`
	Uri             string `json:"uri"               yaml:"uri"               env:"URI"                 validate:"required"`
	AccessKeyId     string `json:"access_key_id"     yaml:"access_key_id"     env:"ACCESS_KEY_ID"       validate:"required"`
	SecretAccessKey string `json:"secret_access_key" yaml:"secret_access_key" env:"SECRET_ACCESS_KEY"   validate:"required"`
}

// Client represents the iaas api client.
type Client struct {
	cfg    *Config
	cli    *ghttp.Client
	tracer gtrace.Tracer
}

// New create a new api client.
func New(ctx context.Context, cfg *Config) *Client {
	_ = ctx
	return &Client{
		cfg:    cfg,
		cli:    ghttp.NewClient(ctx, nil),
		tracer: gtrace.TracerFromContext(ctx),
	}
}

// Build query params; example1: {"action": "DescribeAccessKeys", "access_keys.1": "xxxxxxx"}
func (c *Client) buildQueryParams(params map[string]interface{}) map[string]string {
	queryParams := make(map[string]string)
	for k, v := range params {
		switch temp := v.(type) {
		case []string:
			for i, element := range temp {
				k1 := k + "." + strconv.Itoa(i+1)
				queryParams[k1] = element
			}
		case []int:
			for i, element := range temp {
				k1 := k + "." + strconv.Itoa(i+1)
				queryParams[k1] = strconv.Itoa(element)
			}
		case string:
			// Compatible the upper case zone id in "IaaS"
			queryParams[k] = temp
		case int:
			queryParams[k] = strconv.Itoa(temp)
		default:
			panic(fmt.Sprintf("unsupport data %v for k %s in params", temp, k))
		}
	}
	return queryParams
}

func (c *Client) buildRequestURL(queryParams map[string]string, opts ...Option) string {
	ak := c.cfg.AccessKeyId
	sk := c.cfg.SecretAccessKey

	if len(opts) > 0 {
		var op Operation
		for _, opt := range opts {
			opt(&op)
		}
		if op.accessKeyId != "" {
			ak = op.accessKeyId
		}
		if op.secretAccessKey != "" {
			sk = op.secretAccessKey
		}
	}

	currentTime := time.Now().UTC()
	queryParams["time_stamp"] = currentTime.Format(iso8601)
	queryParams["expires"] = currentTime.Add(time.Second * time.Duration(c.cfg.Timeout)).Format(iso8601)
	queryParams["access_key_id"] = ak
	queryParams["signature_version"] = signatureVersion
	queryParams["signature_method"] = signatureAlgorithm

	// Build signature.
	var queryKeys []string
	for k := range queryParams {
		queryKeys = append(queryKeys, k)
	}
	sort.Strings(queryKeys)

	var queryPairs []string
	for _, k := range queryKeys {
		val := url.QueryEscape(k) + "=" + url.QueryEscape(queryParams[k])
		queryPairs = append(queryPairs, val)
	}
	queryString := strings.Join(queryPairs, "&")

	signToString := strings.Join([]string{
		http.MethodGet,
		c.cfg.Uri,
		queryString,
	}, "\n")

	h := hmac.New(sha256.New, []byte(sk))
	h.Write([]byte(signToString))
	signature := strings.TrimSpace(base64.StdEncoding.EncodeToString(h.Sum(nil)))

	reqURL := fmt.Sprintf(
		"%s://%s:%d%s?%s&signature=%s",
		c.cfg.Protocol,
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.Uri,
		queryString,
		url.QueryEscape(signature),
	)
	return reqURL
}

func (c *Client) sendRequest(ctx context.Context, params map[string]interface{}, respBody ResponseBody, opts ...Option) (err error) {
	var request *http.Request
	var response *http.Response

	lg := glog.FromContext(ctx)

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}
	span := c.tracer.StartSpan(
		"IaasClient",
		opentracing.ChildOf(parentCtx),
		ext.SpanKindRPCClient,
		traceComponentTag,
	)

	// Ensure the request body be close
	defer func() {
		if request != nil && request.Body != nil {
			_ = request.Body.Close()
			request.Body = nil
		}
		if response != nil && response.Body != nil {
			_ = response.Body.Close()
			response.Body = nil
		}

		if err != nil {
			span.LogFields(tracerLog.Error(err))
			lg.Error().Error("send request to iaas error", err).Fire()
		}
		span.Finish()
	}()

	queryParams := c.buildQueryParams(params)
	reqURL := c.buildRequestURL(queryParams, opts...)

	span.LogFields(tracerLog.String("url", reqURL))
	ctx = opentracing.ContextWithSpan(ctx, span)
	if request, err = http.NewRequest(http.MethodGet, reqURL, nil); err != nil {
		return
	}

	//Retry to do sent request
	for i := 0; i < 3; i++ {
		if response, err = c.cli.Send(ctx, request); err != nil {
			time.Sleep(time.Second * time.Duration(i))
			continue
		}
		break
	}
	if err != nil {
		return
	}

	var bodyBytes []byte

	if response.Body != nil && response.ContentLength != 0 {
		if bodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
			lg.Error().Error("read response body error", err).Fire()
			return
		}
		lg.Debug().RawString("response body from iaas", string(bodyBytes)).Fire()
	}

	if response.StatusCode != 200 {
		err = fmt.Errorf("unexpected response status code %d from iaas", response.StatusCode)
		return
	}

	if respBody != nil {
		if err = json.Unmarshal(bodyBytes, respBody); err != nil {
			err = fmt.Errorf("unmarsahl iass response error: %v", err)
			return
		}
		if respBody.ReturnCode() != 0 {
			err = fmt.Errorf("ret_code=%d, message=%s", respBody.ReturnCode(), respBody.ReturnMessage())
			return
		}
	}
	return
}

// DescribeUserById query the user info by specified userId.
func (c *Client) DescribeUserById(ctx context.Context, userId string) (user *User, err error) {
	params := map[string]interface{}{
		"action": "DescribeUsers",
		"users":  []string{userId},
	}
	var body DescribeUsersOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	if len(body.UserSet) == 0 {
		err = ErrUserNotExists
		return
	}
	user = body.UserSet[0]
	return
}

// DescribeUsers get the list of user by giving user ids.
func (c *Client) DescribeUsers(ctx context.Context, input *DescribeUsersInput) (resp *DescribeUsersOutput, err error) {
	params := map[string]interface{}{
		"action": "DescribeUsers",
	}
	if input.Limit != 0 {
		params["limit"] = input.Limit
	}
	if input.Offset != 0 {
		params["offset"] = input.Offset
	}
	if len(input.Users) != 0 {
		params["users"] = input.Users
	}
	if input.Status != "" {
		params["status"] = input.Status
	}
	if input.Email != "" {
		params["email"] = input.Email
	}
	if input.Phone != "" {
		params["phone"] = input.Phone
	}

	var body DescribeUsersOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	resp = &body

	return
}

// DescribeAccessKeysById query the access key info by the specified accessKeyId.
func (c *Client) DescribeAccessKeysById(ctx context.Context, accessKeyId string) (accessKey *AccessKey, err error) {
	params := map[string]interface{}{
		"action":      "DescribeAccessKeys",
		"access_keys": []string{accessKeyId},
		"limit":       1,
		"offset":      0,
	}

	var body DescribeAccessKeysOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	if len(body.AccessKeySet) == 0 {
		err = ErrAccessKeyNotExists
		return
	}
	accessKey = body.AccessKeySet[0]
	return
}

// DescribeAccessKeysByOwner query the pitrix access key info by specified owner.
func (c *Client) DescribeAccessKeysByOwner(ctx context.Context, owner string) (accessKey *AccessKey, err error) {
	params := map[string]interface{}{
		"action":     "DescribeAccessKeys",
		"owner":      owner,
		"controller": "pitrix",
		"limit":      1,
		"offset":     0,
	}

	var body DescribeAccessKeysOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	if len(body.AccessKeySet) == 0 {
		err = ErrAccessKeyNotExists
		return
	}
	accessKey = body.AccessKeySet[0]
	return
}

// DescribeActiveRoutersByOwner only query the active router by giving owner.
func (c *Client) DescribeActiveRoutersByOwner(ctx context.Context, owner string, limit int, offset int) (
	resp *DescribeRoutersOutput, err error) {
	params := map[string]interface{}{
		"action":      "DescribeRouters",
		"routers":     []string{},
		"zone":        c.cfg.Zone,
		"status":      []string{"active"},
		"router_type": []int{1, 0, 2, 3},
		"mode":        0,
		"verbose":     1,
		"owner":       owner,
		"limit":       limit,
		"offset":      offset,
	}

	var body DescribeRoutersOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	resp = &body
	return
}

// DescribeRoutersByOwner query the route info of specified owner.
func (c *Client) DescribeRoutersByOwner(ctx context.Context, owner string, limit int, offset int) (
	resp *DescribeRoutersOutput, err error) {
	params := map[string]interface{}{
		"action":      "DescribeRouters",
		"routers":     []string{},
		"zone":        c.cfg.Zone,
		"status":      []string{"pending", "active", "poweroffed", "suspended"},
		"router_type": []int{1, 0, 2, 3, 99},
		"mode":        0,
		"verbose":     1,
		"owner":       owner,
		"limit":       limit,
		"offset":      offset,
	}

	var body DescribeRoutersOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	resp = &body
	return
}

// DescribeRoutersById query the route info of specified routerId.
func (c *Client) DescribeRoutersById(ctx context.Context, routerId string) (router *Router, err error) {
	params := map[string]interface{}{
		"action":      "DescribeRouters",
		"routers":     []string{routerId},
		"zone":        c.cfg.Zone,
		"status":      []string{"pending", "active", "poweroffed", "suspended"},
		"router_type": []int{99, 1, 0, 2, 3},
		"mode":        0,
		"verbose":     0,
		"limit":       1,
		"offset":      0,
	}

	var body DescribeRoutersOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	if len(body.RouterSet) == 0 {
		err = ErrRouterNotExists
		return
	}
	router = body.RouterSet[0]
	return
}

// DescribeRouterVxnetsById query the router's vxnets by specified routerId.
func (c *Client) DescribeRouterVxnetsById(ctx context.Context, routerId string, limit int, offset int) (resp *DescribeRouterVxnetsOutput, err error) {
	params := map[string]interface{}{
		"action":  "DescribeRouterVxnets",
		"router":  []string{routerId},
		"zone":    c.cfg.Zone,
		"verbose": 0,
		"limit":   limit,
		"offset":  offset,
	}

	var body DescribeRouterVxnetsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	resp = &body
	return
}

// DescribeVxnetById query the vxnet info of specified vxnetId.
func (c *Client) DescribeVxnetById(ctx context.Context, vxnetId string) (vxnet *Vxnet, err error) {
	params := map[string]interface{}{
		"action":          "DescribeVxnets",
		"vxnets":          []string{vxnetId},
		"vxnet_type":      []string{"0", "1"},
		"zone":            c.cfg.Zone,
		"verbose":         1,
		"limit":           1,
		"offset":          0,
		"excluded_vxnets": []string{"vxnet-0", "vxnet-1"},
		//"owner":       "",
	}

	var body DescribeVxnetsOuput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	if len(body.VxnetSet) == 0 {
		err = ErrVXNetNotExists
		return
	}
	vxnet = body.VxnetSet[0]
	return
}

func (c *Client) describeVxnetResources(ctx context.Context, vxnetId string, limit int, offset int, opts ...Option) (
	output *DescribeVxnetResourcesOutput, err error) {
	params := map[string]interface{}{
		"action": "DescribeVxnetResources",
		"vxnet":  vxnetId,
		"zone":   c.cfg.Zone,
		"limit":  limit,
		"offset": offset,
	}

	var body DescribeVxnetResourcesOutput
	if err = c.sendRequest(ctx, params, &body, opts...); err != nil {
		return
	}
	output = &body
	return
}

// DescribeVxnetResources query the vxnet's resources.
// Notice: must access with the owner's assessKey/secretKey.
func (c *Client) DescribeVxnetResources(ctx context.Context, vxnetId string, limit int, offset int, opts ...Option) (
	vxnetResourceSet []*VxnetResource, err error) {

	var output *DescribeVxnetResourcesOutput
	output, err = c.describeVxnetResources(ctx, vxnetId, limit, offset, opts...)
	if err != nil {
		return
	}
	vxnetResourceSet = output.VxnetResourceSet
	return
}

// DescribeAllVxnetResources query the vxnet's all resources.
// Notice: must access with the owner's assessKey/secretKey.
func (c *Client) DescribeAllVxnetResources(ctx context.Context, vxnetId string, opts ...Option) (
	vxnetResourceSet []*VxnetResource, err error) {

	var output *DescribeVxnetResourcesOutput

	setup := 0
	offset := 0
	limit := 100

LOOP:
	for {
		// To avoid dead loop.
		if setup >= 10 {
			err = errors.New("DescribeAllVxnetResources: exception occurred in IaaS")
			return nil, err
		}

		output, err = c.describeVxnetResources(ctx, vxnetId, limit, offset, opts...)
		if err != nil {
			return nil, err
		}
		vxnetResourceSet = append(vxnetResourceSet, output.VxnetResourceSet...)

		if len(vxnetResourceSet) >= output.TotalCount {
			break LOOP
		}
		offset = len(vxnetResourceSet)
		setup++
	}
	return
}

// GetBalance for query the user balance info by specified userId.
func (c *Client) GetBalance(ctx context.Context, userId string) (balanceSet *GetBalanceOutput, err error) {
	params := map[string]interface{}{
		"action": "GetBalance",
		"zone":   c.cfg.Zone,
		"user":   userId,
	}
	var body GetBalanceOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	balanceSet = &body
	return
}

// DescribeJobById is wrapper for DescribeJobs, only query one.
func (c *Client) DescribeJobById(ctx context.Context, jobId string) (jobSet *JobSet, err error) {
	ret, err := c.DescribeJobs(ctx, []string{jobId}, nil)
	if err != nil {
		return
	}
	if len(ret) != 1 {
		err = errors.New("iaas job not exists")
		return
	}
	jobSet = ret[0]
	return
}

func (c *Client) DescribeJobs(ctx context.Context, jobs []string, jobAction []string) (jobSet []*JobSet, err error) {
	params := map[string]interface{}{
		"action":     "DescribeJobs",
		"zone":       c.cfg.Zone,
		"jobs":       jobs,
		"job_action": jobAction,
	}
	var body DescribeJobsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	jobSet = body.JobSet
	return
}

func (c *Client) AllocateVips(ctx context.Context, input *AllocateVipsInput) (jobId string, vips []string, err error) {
	params := map[string]interface{}{
		"action":      "AllocateVips",
		"zone":        c.cfg.Zone,
		"vxnet_id":    input.VxnetId,
		"vip_name":    input.VipName,
		"vip_addrs":   input.VipAddrs,
		"count":       len(input.VipAddrs),
		"target_user": input.TargetUser,
		"vip_range":   input.VipRange,
	}
	var body AllocateVipsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	jobId = body.JobId
	vips = body.Vips
	return
}

func (c *Client) ReleaseVips(ctx context.Context, vips []string) (jobId string, err error) {
	params := map[string]interface{}{
		"action": "ReleaseVips",
		"zone":   c.cfg.Zone,
		"vips":   vips,
	}
	var body ReleaseVipsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	jobId = body.JobId
	return
}

func (c *Client) DescribeVips(ctx context.Context, input *DescribeVipsInput) (output *DescribeVipsOutput, err error) {
	params := map[string]interface{}{
		"action":    "DescribeVips",
		"zone":      c.cfg.Zone,
		"vxnets":    input.Vxnets,
		"limit":     input.Limit,
		"offset":    input.Offset,
		"vips":      input.Vips,
		"vip_addrs": input.VipAddrs,
		//"owner":     input.Owner,
		//"vip_name":  input.VipName,
	}
	if input.Owner != "" {
		params["owner"] = input.Owner
	}
	if input.VipName != "" {
		params["vip_name"] = input.VipName
	}
	var body DescribeVipsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	output = &body
	return
}

func (c *Client) DescribeAllVips(ctx context.Context, input *DescribeVipsInput) (vipSet []*VipSet, err error) {
	var output *DescribeVipsOutput

	setup := 0
	offset := 0
	limit := 100

LOOP:
	for {
		// To avoid dead loop.
		if setup >= 10 {
			err = errors.New("DescribeAllVips: exception occurred in IaaS")
			return nil, err
		}
		input.Limit = limit
		input.Offset = offset

		output, err = c.DescribeVips(ctx, input)
		if err != nil {
			return nil, err
		}
		vipSet = append(vipSet, output.VipSet...)

		if len(vipSet) >= output.TotalCount {
			break LOOP
		}
		offset = len(vipSet)
		setup++
	}
	return
}

func (c *Client) DescribeNotificationLists(ctx context.Context, owner string, nfLists []string, limit int, offset int) (output *DescribeNotificationListsOutput, err error) {
	params := map[string]interface{}{
		"action":             "DescribeNotificationLists",
		"zone":               c.cfg.Zone,
		"owner":              owner,
		"notification_lists": strings.Join(nfLists, ","),
		"limit":              limit,
		"offset":             offset,
	}
	var body DescribeNotificationListsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	output = &body
	return

}
