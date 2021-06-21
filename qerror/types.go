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
	SpaceOwnerCannotRemoved = &Error{
		code:   "SpaceOwnerCannotRemoved",
		status: 403,
		enUS:   "The member of workspace owner cannot be removed.",
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
		enUS:   "not support source type[%s]",
		zhCN:   "不支持的数据源类型[%s]",
	}
	NotSupportEngineType = &Error{
		code:   "NotSupportEngineType",
		status: http.StatusInternalServerError,
		enUS:   "not support engine type[%s]",
		zhCN:   "不支持的引擎类型[%s]",
	}
	InvalidSourceName = &Error{
		code:   "InvalidSourceName",
		status: http.StatusInternalServerError,
		enUS:   "invalid name. can't use '.'",
		zhCN:   "无效名字，不能使用'.'",
	}
	InvalidDimensionSource = &Error{
		code:   "InvalidDimensionSource",
		status: http.StatusInternalServerError,
		enUS:   "dimension just use in relation database",
		zhCN:   "维表只能在关系型数据库使用",
	}
	ConnectSourceFailed = &Error{
		code:   "ConnectSourceFailed",
		status: http.StatusInternalServerError,
		enUS:   "this source can't connect",
		zhCN:   "数据源无法连接",
	}
)

// udfmanager
var (
	InvalidUdfName = &Error{
		code:   "InvalidUdfName",
		status: http.StatusInternalServerError,
		enUS:   "invalid name. can't use '.'",
		zhCN:   "无效名字，不能使用'.'",
	}
	NotSupportUdfType = &Error{
		code:   "NotSupportUdfType",
		status: http.StatusInternalServerError,
		enUS:   "not support udf type[%s]",
		zhCN:   "不支持的自定义函数类型[%s]",
	}
)

// enginemanager error
var (
	EngineNameAlreadyExist = &Error{
		code:   "EngineNameAlreadyExist",
		status: 409,
		enUS:   "The engine name[%s] has been used.",
		zhCN:   "计算引擎名称[%s]已占用.",
	}

	EngineNotExist = &Error{
		code:   "EngineNotExist",
		status: 404,
		enUS:   "The engine[%s] not exist.",
		zhCN:   "计算引擎[%s]不存在.",
	}
)

// notification error
var (
	SendingNotificationFailed = &Error{
		code:   "SendingNotificationFailed",
		status: 500,
		enUS:   "failed to send notification post to user [%s]",
		zhCN:   "向用户[%s]发送通知失败.",
	}

	NotificationChannelConfigInvalid = &Error{
		code:   "NotificationChannelConfigInvalid",
		status: 500,
		enUS:   "notification channel config not valid",
		zhCN:   "通知渠道配置不合法",
	}
)
