package qerror

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
	WorkspaceNotExists = &Error{
		code:   "WorkspaceNotExists",
		status: 404,
		enUS:   "The workspace [%s] does not exists.",
		zhCN:   "工作空间[%s]不存在.",
	}
	WorkspaceAlreadyExists = &Error{
		code:   "WorkspaceAlreadyExists",
		status: 409,
		enUS:   "The workspace name [%s] has been used.",
		zhCN:   "工作空间名称[%s]已被使用.",
	}
	WorkspaceNotActive = &Error{
		code:   "WorkspaceNotActive",
		status: 403,
		enUS:   "The workspace [%s] does not active.",
		zhCN:   "工作空间[%s]已被禁用",
	}
	InvalidWorkspaceName = &Error{
		code:   "InvalidWorkspaceName",
		status: 400,
		enUS:   "The workspace name is invalid, accepts 0~9、a~z、_ and can't begin or end with _.",
		zhCN:   "工作空间名称不符合要求, 只允许数字,小写字母和下划线, 并且不能以下划线开头或者结尾.",
	}
)

// workflow error
var (
	WorkflowAlreadyExists = &Error{
		code:   "WorkflowAlreadyExists",
		status: 409,
		enUS:   "The workflow name [%s] has been used.",
		zhCN:   "工作流程名称[%s]已被使用.",
	}
	WorkflowNotExists = &Error{
		code:   "WorkflowNotExists",
		status: 404,
		enUS:   "The workflow [%s] does not exists.",
		zhCN:   "工作流程[%s]不存在.",
	}
	CrontabNotExists = &Error{
		code:   "CrontabNotExists",
		status: 404,
		enUS:   "The workflow [%s] not set crontab.",
		zhCN:   "工作流程[%s]未设置定时任务",
	}
	InvalidWorkflowName = &Error{
		code:   "InvalidWorkflowName",
		status: 400,
		enUS:   "The workflow name is invalid, accepts 0~9、a~z、_ and can't begin or end with _.",
		zhCN:   "工作流程名称不符合要求, 只允许数字,小写字母和下划线, 并且不能以下划线开头或者结尾.",
	}
)

// node error
var (
	NodeNotExists = &Error{
		code:   "NodeNotExists",
		status: 404,
		enUS:   "The node [%s] does not exists.",
		zhCN:   "任务节点[%s]不存在.",
	}
	NodeAlreadyExists = &Error{
		code:   "NodeAlreadyExists",
		status: 409,
		enUS:   "The node name [%s] has been used.",
		zhCN:   "任务节点名称[%s]已被使用.",
	}
	InvalidNodeName = &Error{
		code:   "InvalidNodeName",
		status: 400,
		enUS:   "The node name is invalid, accepts 0~9、a~z、_ and can't begin or end with _.",
		zhCN:   "节点名称不符合要求, 只允许数字,小写字母和下划线, 并且不能以下划线开头或者结尾.",
	}
)
