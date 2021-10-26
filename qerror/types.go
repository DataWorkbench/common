package qerror

import (
	"net/http"
)

// general error
var (
	Internal = &Error{
		code:   "InternalError",
		status: 500,
		enUS:   "We encountered an internal error, Please try again.",
		zhCN:   "内部错误, 请稍后重试.",
	}

	MethodNotAllowed = &Error{
		code:   "MethodNotAllowed",
		status: 405,
		enUS:   "The specified method is not allowed against this resource.",
		zhCN:   "请求被拒绝.",
	}

	NotAcceptable = &Error{
		code:   "NotAcceptable",
		status: 406,
		enUS:   "The request is not acceptable.",
		zhCN:   "请求不可达.",
	}

	// PermissionDenied render message when use have no permission to accessing resource
	PermissionDenied = &Error{
		code:   "PermissionDenied",
		status: 403,
		enUS:   "You don't have enough permission to accomplish this request",
		zhCN:   "权限拒绝",
	}

	// InvalidJSON render message if JSON format error.
	InvalidJSON = &Error{
		code:   "InvalidJSON",
		status: 400,
		enUS:   "The request body is invalid as JSON format.",
		zhCN:   "JSON 格式错误",
	}

	UnknownField = &Error{
		code:   "UnknownField",
		status: 400,
		enUS:   "Found unknown field [%s] in you request body",
		zhCN:   "请求体中包含不支持字段[%s]",
	}

	InvalidRequest = &Error{
		code:   "InvalidRequest",
		status: 400,
		enUS:   "The requests you provided is invalid.",
		zhCN:   "无效的请求",
	}

	// NotActive render message if accessing resource is inactive and operation not allows.
	ResourceNotActive = &Error{
		code:   "ResourceNotActive",
		status: 403,
		enUS:   "The resource you are accessing does not active.",
		zhCN:   "资源已被禁用, 请先启用.",
	}

	// NotExists render message if accessing resource not exists.
	ResourceNotExists = &Error{
		code:   "ResourceNotExists",
		status: 404,
		enUS:   "The resource you are accessing does not exist.",
		zhCN:   "资源不存在",
	}

	// AlreadyExists render message if creating resource already exists.
	ResourceAlreadyExists = &Error{
		code:   "ResourceAlreadyExists",
		status: 409,
		enUS:   "The resource you are creating already exists.",
		zhCN:   "资源已存在",
	}

	ResourceIsUsing = &Error{
		code:   "ResourceIsUsing",
		status: http.StatusProcessing,
		enUS:   "The resource is using.",
		zhCN:   "资源正在使用",
	}
)

// parameters error
var (
	// InvalidParams render message if request parameter verification failed
	InvalidParams = &Error{
		code:   "InvalidParams",
		status: 400,
		enUS:   "The parameters [%s] you provided is invalid.",
		zhCN:   "参数[%s]错误",
	}
	ParamsIsEmpty = &Error{
		code:   "ParamsIsEmpty",
		status: 400,
		enUS:   "The parameters [%s] is required.",
		zhCN:   "参数[%s]不能为空.",
	}
	InvalidParamsLength = &Error{
		code:   "InvalidParamsLength",
		status: 400,
		enUS:   "The length of parameters [%s] must be [%s %s] and you provided parameters length [%d].",
		zhCN:   "参数[%s]的长度必须 [%s %s], 请求参数的长度 [%d].",
	}
	InvalidParamsValue = &Error{
		code:   "InvalidParamsValue",
		status: 400,
		enUS:   "The value of [%s] must be [%s %s] and you provides value is [%v].",
		zhCN:   "参数[%s]必须 [%s %s], 请求参数值 [%v].",
	}
	DependenceParamNotExist = &Error{
		code:   "DependenceParamNotExist",
		status: 400,
		enUS:   "The param[s] [%s] is[are] dependence of param [%s] not exist.",
		zhCN:   "参数[%s]不存在，被参数[%s]依赖.",
	}
)

// Error for Global API
var (
	RegionNotSpecified = &Error{
		code:   "RegionNotSpecified",
		status: 400,
		enUS:   "A valid region id must be specified in you request path",
		zhCN:   "无效的请求, 未指定 regionId",
	}
	RegionNotExists = &Error{
		code:   "RegionNotExists",
		status: 404,
		enUS:   "The region [%s] you access not exists.",
		zhCN:   "访问的区域[%s]不存在.",
	}
	RegionAccessDenied = &Error{
		code:   "RegionAccessDenied",
		status: 403,
		enUS:   "The user [%s] is not allowed to access region [%s].",
		zhCN:   "用户 [%s] 没有访问区域 [%s] 的权限",
	}
)

// workspace error
var (
	SpaceNotExists = &Error{
		code:   "SpaceNotExists",
		status: 404,
		enUS:   "The workspace [%s] does not exists.",
		zhCN:   "工作空间[%s]不存在.",
	}
	SpaceAlreadyExists = &Error{
		code:   "SpaceAlreadyExists",
		status: 409,
		enUS:   "The workspace name [%s] has been used.",
		zhCN:   "工作空间名称[%s]已被使用.",
	}
	SpaceNotActive = &Error{
		code:   "SpaceNotActive",
		status: 403,
		enUS:   "The workspace [%s] does not active.",
		zhCN:   "工作空间[%s]已被禁用",
	}
	InvalidSpaceName = &Error{
		code:   "InvalidWorkspaceName",
		status: 400,
		enUS:   "The workspace name is invalid, accepts 0~9、a~z、_ and can't begin or end with _.",
		zhCN:   "工作空间名称不符合要求, 只允许数字,小写字母和下划线, 并且不能以下划线开头或者结尾.",
	}
)

// member error
var (
	MemberNotExists = &Error{
		code:   "MemberNotExists",
		status: 404,
		enUS:   "The member [%s] does not exists.",
		zhCN:   "成员[%s]不存在.",
	}
	MemberAlreadyExists = &Error{
		code:   "MemberAlreadyExists",
		status: 409,
		enUS:   "The member [%s] has been exists.",
		zhCN:   "空间成员[%s]已存在.",
	}
	SpaceOwnerCannotDelete = &Error{
		code:   "SpaceOwnerCannotDelete",
		status: 403,
		enUS:   "The member of workspace owner cannot be deleted.",
		zhCN:   "空间所有者不允许删除.",
	}
	SpaceOwnerCannotUpdated = &Error{
		code:   "SpaceOwnerCannotUpdated",
		status: 403,
		enUS:   "The member of workspace owner cannot be updated.",
		zhCN:   "空间所有者不允许更新",
	}
)

// workflow error
var (
	FlowAlreadyExists = &Error{
		code:   "FlowAlreadyExists",
		status: 409,
		enUS:   "The workflow name [%s] has been used.",
		zhCN:   "工作流程名称[%s]已被使用.",
	}
	FlowNotExists = &Error{
		code:   "FlowNotExists",
		status: 404,
		enUS:   "The workflow [%s] does not exists.",
		zhCN:   "工作流程[%s]不存在.",
	}
	ScheduleNotSet = &Error{
		code:   "ScheduleNotSet",
		status: 400,
		enUS:   "The workflow [%s] not set schedule properties.",
		zhCN:   "工作流程[%s]未设置调度属性",
	}
	NodeNotSet = &Error{
		code:   "NodeNotSet",
		status: 400,
		enUS:   "The workflow [%s] not set node task.",
		zhCN:   "工作流程[%s]未设置节点任务",
	}
	EnvNotSet = &Error{
		code:   "EnvNotSet",
		status: 400,
		enUS:   "The workflow [%s] not set environmental parameters ",
		zhCN:   "工作流程[%s]未设置环境参数",
	}
	InvalidFlowName = &Error{
		code:   "InvalidFlowName",
		status: 400,
		enUS:   "The workflow name is invalid, accepts 0~9、a~z、_ and can't begin or end with _.",
		zhCN:   "工作流程名称不符合要求, 只允许数字,小写字母和下划线, 并且不能以下划线开头或者结尾.",
	}
)

// instance error
var (
	InstanceNotExists = &Error{
		code:   "InstanceNotExists",
		status: 404,
		enUS:   "The instance [%s] does not exists.",
		zhCN:   "任务实例[%s]不存在.",
	}
	//NodeAlreadyExists = &Error{
	//	code:   "NodeAlreadyExists",
	//	status: 409,
	//	enUS:   "The node name [%s] has been used.",
	//	zhCN:   "任务节点名称[%s]已被使用.",
	//}
	//InvalidNodeName = &Error{
	//	code:   "InvalidNodeName",
	//	status: 400,
	//	enUS:   "The node name is invalid, accepts 0~9、a~z、_ and can't begin or end with _.",
	//	zhCN:   "节点名称不符合要求, 只允许数字,小写字母和下划线, 并且不能以下划线开头或者结尾.",
	//}
)

// sourcemanager
var (
	NotSupportSourceType = &Error{
		code:   "NotSupportSourceType",
		status: http.StatusInternalServerError,
		enUS:   "not support source type[%d]",
		zhCN:   "不支持的数据源类型[%d]",
	}
	InvalidSourceName = &Error{
		code:   "InvalidSourceName",
		status: http.StatusInternalServerError,
		enUS:   "invalid name. can't use '.'",
		zhCN:   "无效名字，不能使用'.'",
	}
	SourceIsDisable = &Error{
		code:   "SourceIsDisable",
		status: http.StatusInternalServerError,
		enUS:   "source is disable",
		zhCN:   "数据源是禁用状态",
	}
	ConnectSourceFailed = &Error{
		code:   "ConnectSourceFailed",
		status: http.StatusInternalServerError,
		enUS:   "this source can't connect[%d]",
		zhCN:   "数据源无法连接[%d]",
	}
)

// udfmanager
var (
	InvalidUDFName = &Error{
		code:   "InvalidUDFName",
		status: http.StatusInternalServerError,
		enUS:   "invalid name. can't use '.'",
		zhCN:   "无效名字，不能使用'.'",
	}
	NotSupportUDFType = &Error{
		code:   "NotSupportUDFType",
		status: http.StatusInternalServerError,
		enUS:   "not support UDF type[%d]",
		zhCN:   "不支持的自定义函数类型[%d]",
	}
)

// enginemanager error
var (
	EngineNotExist = &Error{
		code:   "EngineNotExist",
		status: 404,
		enUS:   "The engine[%s] not exist.",
		zhCN:   "计算引擎[%s]不存在.",
	}

	EngineNameInUse = &Error{
		code:   "EngineNameInUse",
		status: 400,
		enUS:   "The engine name[%s] is in use.",
		zhCN:   "计算引擎名称[%s]已存在.",
	}

	EngineNameError = &Error{
		code:   "EngineNameError",
		status: 400,
		enUS:   "The name[%s] is not support.",
		zhCN:   "命名[%s]错误.",
	}

	EngineInTransaction = &Error{
		code:   "EngineInTransaction",
		status: 403,
		enUS:   "The engine[%s] is in transaction[%s].",
		zhCN:   "计算引擎[%s]正在[%s]中.",
	}

	EngineCreateTimeout = &Error{
		code:   "EngineCreateTimeout",
		status: 500,
		enUS:   "Create Engine[%s] timeout.",
		zhCN:   "创建引擎[%s]超时.",
	}
)

// notifier error
var (
	SendingNotificationFailed = &Error{
		code:   "SendingNotificationFailed",
		status: 500,
		enUS:   "failed to send notifier post to user [%s]",
		zhCN:   "向用户[%s]发送通知失败.",
	}

	NotificationChannelConfigInvalid = &Error{
		code:   "NotificationChannelConfigInvalid",
		status: 500,
		enUS:   "notifier channel config not valid",
		zhCN:   "通知渠道配置不合法",
	}
)

// account error
var (
	AccessKeyNotExists = &Error{
		code:   "AccessKeyNotExist",
		status: 404,
		enUS:   "The access key[%s] not exist.",
		zhCN:   "Access key[%s]不存在。",
	}
	UserNotExists = &Error{
		code:   "UserNotExist",
		status: 404,
		enUS:   "The user[%s] not exist.",
		zhCN:   "User[%s]不存在。",
	}
	ValidateSignatureFailed = &Error{
		code:   "ValidateSignatureFailed",
		status: 400,
		enUS:   "Validate signature failed, the signature[%s] not match [%s].",
		zhCN:   "签名验证失败, 签名[%s]与[%s]不匹配。",
	}
)

// logmanager error
var (
	LogFileNotExists = &Error{
		code:   "LogFileNotExists",
		status: 405,
		enUS:   "Log file [%s] not exists",
		zhCN:   "日志文件 [%s] 不存在",
	}

	RequestForFlinkFailed = &Error{
		code:   "RequestForFlinkFailed",
		status: 406,
		enUS:   "fail to request flink api [%s]",
		zhCN:   "调用Flink Web Api [%s] 失败",
	}
)
