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
	requestPath        = "/iaas/"
)

// Config represents the iaas api config.
type Config struct {
	Zone            string `json:"zone"              yaml:"zone"              env:"ZONE"                validate:"required"`
	Host            string `json:"host"              yaml:"host"              env:"HOST"                validate:"required"`
	Port            int    `json:"port"              yaml:"port"              env:"PORT"                validate:"required"`
	Protocol        string `json:"protocol"          yaml:"protocol"          env:"PROTOCOL"            validate:"required"`
	Timeout         int    `json:"timeout"           yaml:"timeout"           env:"TIMEOUT,default=600" validate:"required"`
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
		requestPath,
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
		requestPath,
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
	if response.StatusCode != 200 {
		err = fmt.Errorf("unexpected response status code %d from iaas", response.StatusCode)
		return
	}

	if respBody != nil {
		var bodyBytes []byte
		if bodyBytes, err = ioutil.ReadAll(response.Body); err != nil {
			return
		}

		lg.Debug().RawString("response body from iaas", string(bodyBytes)).Fire()

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

// DescribeUsersById query the user info by specified userId.
func (c *Client) DescribeUsersById(ctx context.Context, userId string) (user *User, err error) {
	params := map[string]interface{}{
		"action": "DescribeUsers",
		"users":  []string{userId},
	}
	var body DescribeUsersOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	if len(body.UserSet) == 0 {
		err = errors.New("user_not_exists")
		return
	}
	user = body.UserSet[0]
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
		err = errors.New("access_key_not_exists")
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
		err = errors.New("access_key_not_exists")
		return
	}
	accessKey = body.AccessKeySet[0]
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
		err = errors.New("router_key_not_exists")
		return
	}
	router = body.RouterSet[0]
	return
}

// DescribeRouterVxnetsById query the router's vxnets by specified routerId.
func (c *Client) DescribeRouterVxnetsById(ctx context.Context, routerId string) (routerVxnetSet []*RouterVxnet, err error) {
	params := map[string]interface{}{
		"action":  "DescribeRouterVxnets",
		"router":  []string{routerId},
		"zone":    c.cfg.Zone,
		"verbose": 0,
		"limit":   100,
		"offset":  0,
	}

	var body DescribeRouterVxnetsOutput
	if err = c.sendRequest(ctx, params, &body); err != nil {
		return
	}
	routerVxnetSet = body.RouterVxnetSet
	return
}

// DescribeVxnetsById query the vxnet info of specified vxnetId.
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
		err = errors.New("vxnet_not_exists")
		return
	}
	vxnet = body.VxnetSet[0]
	return
}

// DescribeVxnetResources query the vxnet's resources.
// Notice: must access with the owner's assessKey/secretKey.
func (c *Client) DescribeVxnetResources(ctx context.Context, vxnetId string, limit int, offset int, opts ...Option) (
	vxnetResourceSet []*VxnetResource, err error) {

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
	vxnetResourceSet = body.VxnetResourceSet
	return
}
