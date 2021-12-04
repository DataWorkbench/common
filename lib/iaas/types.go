package iaas

type ResponseBody interface {
	ReturnCode() int
	ReturnMessage() string
}

// User represents the user info.
type User struct {
	UserId     string   `json:"user_id"`
	Name       string   `json:"user_name"`
	Email      string   `json:"email"`
	RootUserId string   `json:"root_user_id"`
	Role       string   `json:"role"`
	Status     string   `json:"status"`
	Privilege  int      `json:"privilege"`
	Zones      []string `json:"zones"`
	Regions    []string `json:"regions"`
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
type AccessKey struct {
	AccessKeyId     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	Owner           string `json:"owner"`
	RootUserId      string `json:"root_user_id"`
}

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
type Router struct {
	RouterId   string `json:"router_id"`
	RouterName string `json:"router_name"`
	RouterType int    `json:"router_type"`
	Owner      string `json:"owner"`
	Status     string `json:"status"`
	BaseVxnet  string `json:"base_vxnet"`
	VPCNetwork string `json:"vpc_network"`
	PrivateIP  string `json:"private_ip"`
	VPCId      string `json:"vpc_id"`

	// Field for DescribeVxnets
	ManagerIp       string `json:"manager_ip"`
	IpNetwork       string `json:"ip_network"`
	DynIpStart      string `json:"dyn_ip_start"`
	DynIpEnd        string `json:"dyn_ip_end"`
	BorderPrivateIp string `json:"border_private_ip"`
}

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
type RouterVxnet struct {
	RouterId        string `json:"router_id"`
	VxnetId         string `json:"vxnet_id"`
	VxnetName       string `json:"vxnet_name"`
	DynIpStart      string `json:"dyn_ip_start"`
	DynIpEnd        string `json:"dyn_ip_end"`
	DynIpv6Start    string `json:"dyn_ipv_6_start"`
	DynIpv6End      string `json:"dyn_ipv_6_end"`
	Owner           string `json:"owner"`
	BorderPrivateIp string `json:"border_private_ip"`
	ManagerIp       string `json:"manager_ip"`
	BorderId        string `json:"border_id"`
	IpNetwork       string `json:"ip_network"`
	Ipv6Network     string `json:"ipv_6_network"`
	Mode            int    `json:"mode"`
	VpcId           string `json:"vpc_id"`
}

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
type Vxnet struct {
	VxnetId     string  `json:"vxnet_id"`
	VxnetName   string  `json:"vxnet_name"`
	VxnetType   int     `json:"vxnet_type"`
	Owner       string  `json:"owner"`
	TunnelType  string  `json:"tunnel_type"`
	VpcRouterId string  `json:"vpc_router_id"`
	Router      *Router `json:"router"`
}

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
	ResourceType int    `json:"resource_type"`
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
