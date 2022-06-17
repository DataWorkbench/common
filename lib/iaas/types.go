package iaas

import (
	"encoding/json"

	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel/pbiaas"
)

type ResponseBody interface {
	ReturnCode() int
	ReturnMessage() string
}

// User represents the user info.
type User = pbmodel.User

// DescribeUsersInput is type request parameters for action "DescribeUsers"
type DescribeUsersInput struct {
	// The list of user id.
	Users []string

	Limit int

	Offset int

	Status string

	Email string

	Phone string
}

// DescribeUsersOutput is type response body for action "DescribeUsers"
type DescribeUsersOutput struct {
	RetCode    int     `json:"ret_code"`
	Message    string  `json:"message"`
	TotalCount int     `json:"total_count"`
	UserSet    []*User `json:"user_set"`
}

func (b *DescribeUsersOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeUsersOutput) ReturnMessage() string {
	return b.Message
}

// AccessKey represents the access key info.
type AccessKey = pbmodel.AccessKey

//type AccessKey struct {
//	AccessKeyId     string `json:"access_key_id"`
//	SecretAccessKey string `json:"secret_access_key"`
//	Owner           string `json:"owner"`
//	RootUserId      string `json:"root_user_id"`
//}

// DescribeAccessKeysOutput is type response body for action "DescribeAccessKeys"
type DescribeAccessKeysOutput struct {
	RetCode      int          `json:"ret_code"`
	Message      string       `json:"message"`
	TotalCount   int          `json:"total_count"`
	AccessKeySet []*AccessKey `json:"access_key_set"`
}

func (b *DescribeAccessKeysOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeAccessKeysOutput) ReturnMessage() string {
	return b.Message
}

// Router represents the router info.
type Router = pbiaas.Router

//type Router struct {
//	RouterId   string `json:"router_id"`
//	RouterName string `json:"router_name"`
//	RouterType int    `json:"router_type"`
//	Owner      string `json:"owner"`
//	Status     string `json:"status"`
//	BaseVxnet  string `json:"base_vxnet"`
//	VPCNetwork string `json:"vpc_network"`
//	PrivateIP  string `json:"private_ip"`
//	VPCId      string `json:"vpc_id"`
//
//	// Field for DescribeVxnets
//	ManagerIp       string `json:"manager_ip"`
//	IpNetwork       string `json:"ip_network"`
//	DynIpStart      string `json:"dyn_ip_start"`
//	DynIpEnd        string `json:"dyn_ip_end"`
//	BorderPrivateIp string `json:"border_private_ip"`
//}

// DescribeRoutersOutput is type response body for action "DescribeRouters"
type DescribeRoutersOutput struct {
	RetCode    int       `json:"ret_code"`
	Message    string    `json:"message"`
	TotalCount int       `json:"total_count"`
	RouterSet  []*Router `json:"router_set"`
}

func (b *DescribeRoutersOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeRoutersOutput) ReturnMessage() string {
	return b.Message
}

// RouterVxnet represetns the RouterVxnet info.
type RouterVxnet = pbiaas.RouterVxnet

//type RouterVxnet struct {
//	RouterId        string `json:"router_id"`
//	VxnetId         string `json:"vxnet_id"`
//	VxnetName       string `json:"vxnet_name"`
//	DynIpStart      string `json:"dyn_ip_start"`
//	DynIpEnd        string `json:"dyn_ip_end"`
//	DynIpv6Start    string `json:"dyn_ipv_6_start"`
//	DynIpv6End      string `json:"dyn_ipv_6_end"`
//	Owner           string `json:"owner"`
//	BorderPrivateIp string `json:"border_private_ip"`
//	ManagerIp       string `json:"manager_ip"`
//	BorderId        string `json:"border_id"`
//	IpNetwork       string `json:"ip_network"`
//	Ipv6Network     string `json:"ipv_6_network"`
//	Mode            int    `json:"mode"`
//	VpcId           string `json:"vpc_id"`
//}

// DescribeRouterVxnetsOutput is the type of response body of action "DescribeRouterVxnets"
type DescribeRouterVxnetsOutput struct {
	RetCode        int            `json:"ret_code"`
	Message        string         `json:"message"`
	TotalCount     int            `json:"total_count"`
	RouterVxnetSet []*RouterVxnet `json:"router_vxnet_set"`
}

func (b *DescribeRouterVxnetsOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeRouterVxnetsOutput) ReturnMessage() string {
	return b.Message
}

// Vxnet represents the vxnet info.
type Vxnet = pbiaas.VXNet

//type Vxnet struct {
//	VxnetId     string  `json:"vxnet_id"`
//	VxnetName   string  `json:"vxnet_name"`
//	VxnetType   int     `json:"vxnet_type"`
//	Owner       string  `json:"owner"`
//	TunnelType  string  `json:"tunnel_type"`
//	VpcRouterId string  `json:"vpc_router_id"`
//	Router      *Router `json:"router"`
//}

// DescribeVxnetsOuput is the type of response body of action "DescribeVxnets"
type DescribeVxnetsOuput struct {
	RetCode    int      `json:"ret_code"`
	Message    string   `json:"message"`
	TotalCount int      `json:"total_count"`
	VxnetSet   []*Vxnet `json:"vxnet_set"`
}

func (b *DescribeVxnetsOuput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeVxnetsOuput) ReturnMessage() string {
	return b.Message
}

// VxnetResource represents the resource info in vxnet.
type VxnetResource struct {
	VxnetId      string `json:"vxnet_id"`
	ResourceId   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Status       string `json:"status"`
	PrivateIp    string `json:"private_ip"`
	Owner        string `json:"owner"`
}

// DescribeVxnetResourcesOutput is the type of response body of action "DescribeVxnetResources"
type DescribeVxnetResourcesOutput struct {
	RetCode          int              `json:"ret_code"`
	Message          string           `json:"message"`
	TotalCount       int              `json:"total_count"`
	VxnetResourceSet []*VxnetResource `json:"vxnet_resource_set"`
}

func (b *DescribeVxnetResourcesOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeVxnetResourcesOutput) ReturnMessage() string {
	return b.Message
}

type CouponsConditions struct {
	Zones         []string `json:"zones"`
	ResourceTypes []string `json:"resource_types"`
	Apps          []string `json:"apps"`
}

type Coupons struct {
	// Status: activated, ?
	Status    string `json:"status"`
	Balance   string `json:"balance"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`

	// "[{\"zones\":[\"all\"],\"resource_types\":[\"all\", \"cluster_app_service\"],\"apps\":[{\"app_id\": \"app-00r26u27\"}]}]"
	Conditions string `json:"conditions"`

	// Thus field may be useful.
	UserId     string `json:"user_id"`
	RootUserId string `json:"root_user_id"`
	UsageMode  string `json:"usage_mode"`
	CouponId   string `json:"coupon_id"`

	// Thus field is useless.
	//
	//Category       string      `json:"category"`
	//ParentCouponId string      `json:"parent_coupon_id"`
	//SubCategory    string      `json:"sub_category"`
	//ResourceId     string      `json:"resource_id"`
	//CouponTypeId   interface{} `json:"coupon_type_id"`
	//Value          string      `json:"value"`
	//Remarks        interface{} `json:"remarks"`
	//ConsoleId      string      `json:"console_id"`
	//Dispatcher     interface{} `json:"dispatcher"`
	//CreateTime     string      `json:"create_time"`
	//StatusTime     interface{} `json:"status_time"`
	//UpdateTime     string      `json:"update_time"`
}

func (c *Coupons) ParseConditions() (output []*CouponsConditions, err error) {
	if c.Conditions == "" {
		return
	}
	output = make([]*CouponsConditions, 0)
	err = json.Unmarshal([]byte(c.Conditions), &output)
	if err != nil {
		return
	}
	return
}

const BillingPaidMode = "prepaid"

// GetBalanceOutput is the type of response body of action "GetBalance"
type GetBalanceOutput struct {
	RetCode int    `json:"ret_code"`
	Message string `json:"message"`

	PaidMode      string     `json:"paid_mode"` // prepaid, ?
	Bonus         string     `json:"bonus"`
	SharedBonus   string     `json:"shared_bonus"`
	Balance       string     `json:"balance"`
	SharedBalance string     `json:"shared_balance"`
	Coupons       []*Coupons `json:"coupons"`
	SharedCoupons []*Coupons `json:"shared_coupons"`

	// Thus field may be useful.
	UserId         string `json:"user_id"`
	RootUserId     string `json:"root_user_id"`
	UserType       int    `json:"user_type"`
	TotalSum       string `json:"total_sum"`
	TotalSharedSum string `json:"total_shared_sum"`

	// Thus field is useless.
	//
	//SharedPaidMode string        `json:"shared_paid_mode"`
	//IncomeHkd      string        `json:"income_hkd"`
	//IncomeUsd      string        `json:"income_usd"`
	//Preference     int           `json:"preference"`
	//IncomeCny      string        `json:"income_cny"`
}

func (b *GetBalanceOutput) ReturnCode() int {
	return b.RetCode
}

func (b *GetBalanceOutput) ReturnMessage() string {
	return b.Message
}

const (
	JobSetStatusSuccessful = "successful"
)

type JobSet struct {
	Status    string `json:"status"`
	JobAction string `json:"job_action"`
}

type DescribeJobsOutput struct {
	RetCode    int       `json:"ret_code"`
	Message    string    `json:"message"`
	TotalCount int       `json:"total_count"`
	JobSet     []*JobSet `json:"job_set"`
}

func (b *DescribeJobsOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeJobsOutput) ReturnMessage() string {
	return b.Message
}

type DescribeVipsInput struct {
	Limit    int      `json:"limit"`
	Offset   int      `json:"offset"`
	Vxnets   []string `json:"vxnets"`    // list of vxnet id.
	Vips     []string `json:"vips"`      // list id of vip. e.g []string{"vip-xxx", "vip-xxx}
	VipAddrs []string `json:"vip_addrs"` // ip address of vip. e.g: []string{"127.0.0.1", "127.0.0.1"}
	VipName  string   `json:"vip_name"`
	Owner    string   `json:"owner"`
}

type AllocateVipsInput struct {
	VxnetId    string   `json:"vxnet_id"`    // required.
	VipName    string   `json:"vip_name"`    // required
	TargetUser string   `json:"target_user"` // required. the user who will own this vip.
	VipAddrs   []string `json:"vip_addrs"`   // e.g: []string{"172.20.0.105", "172.20.0.106", "172.20.0.107"},
	VipRange   string   `json:"vip_range"`   // e.g: 172.20.0.105-172.20.0.108, contains 105 and 108. mutex with VipAddrs.
}

type VipSet struct {
	InstanceNaem string      `json:"instance_naem"`
	EipId        string      `json:"eip_id"`
	VxnetId      string      `json:"vxnet_id"`
	VipId        string      `json:"vip_id"`
	RootUserId   string      `json:"root_user_id"`
	ConsoleId    string      `json:"console_id"`
	InstanceId   string      `json:"instance_id"`
	Controller   string      `json:"controller"`
	CreateTime   string      `json:"create_time"`
	VipAddr      string      `json:"vip_addr"`
	NeedSg       int         `json:"need_sg"`
	Owner        string      `json:"owner"`
	NicId        string      `json:"nic_id"`
	VipName      string      `json:"vip_name"`
	Description  interface{} `json:"description"`
}

type DescribeVipsOutput struct {
	RetCode    int       `json:"ret_code"`
	Message    string    `json:"message"`
	TotalCount int       `json:"total_count"`
	VipSet     []*VipSet `json:"vip_set"`
}

func (b *DescribeVipsOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeVipsOutput) ReturnMessage() string {
	return b.Message
}

type AllocateVipsOutput struct {
	RetCode int      `json:"ret_code"`
	Message string   `json:"message"`
	Vips    []string `json:"vips"` // the list of vip id.
	JobId   string   `json:"job_id"`
}

func (b *AllocateVipsOutput) ReturnCode() int {
	return b.RetCode
}

func (b *AllocateVipsOutput) ReturnMessage() string {
	return b.Message
}

type ReleaseVipsOutput struct {
	RetCode int    `json:"ret_code"`
	Message string `json:"message"`
	JobId   string `json:"job_id"`
}

func (b *ReleaseVipsOutput) ReturnCode() int {
	return b.RetCode
}

func (b *ReleaseVipsOutput) ReturnMessage() string {
	return b.Message
}

type NotificationListSet = pbmodel.NotificationList

type DescribeNotificationListsOutput struct {
	RetCode             int                    `json:"ret_code"`
	Message             string                 `json:"message"`
	TotalCount          int                    `json:"total_count"`
	NotificationListSet []*NotificationListSet `json:"notification_list_set"`
}

func (b *DescribeNotificationListsOutput) ReturnCode() int {
	return b.RetCode
}

func (b *DescribeNotificationListsOutput) ReturnMessage() string {
	return b.Message
}
