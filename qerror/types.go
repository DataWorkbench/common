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
	InvalidName = &Error{
		code:   "InvalidName",
		status: 400,
		enUS:   "The name is invalid, accepts 0~9、a~z、 A~Z、_ and can't begin or end with _.",
		zhCN:   "名称不符合要求, 只允许数字,大小写字母和下划线, 并且不能以下划线开头或者结尾.",
	}

	ParameterValidationError = &Error{
		code:   "ParameterValidationError",
		status: 400,
		enUS:   "%s",
		zhCN:   "%s",
	}
)

// Error for quota.
var (
	QuotaInsufficientWorkspaceLimit = &Error{
		code:   "QuotaInsufficientWorkspaceLimit",
		status: 403,
		enUS:   "Limit exceeded of workspace number quota. quota is [%d] and already use [%d].",
		zhCN:   "工作空间个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientStreamJobLimit = &Error{
		code:   "QuotaInsufficientStreamJobLimit",
		status: 403,
		enUS:   "Limit exceeded of steam job number quota. quota is [%d] and already use [%d].",
		zhCN:   "实时计算作业个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientDataSourceLimit = &Error{
		code:   "QuotaInsufficientDataSourceLimit",
		status: 403,
		enUS:   "Limit exceeded of data source number quota. quota is [%d] and already use [%d].",
		zhCN:   "数据源个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientDataUDFLimit = &Error{
		code:   "QuotaInsufficientDataUDFLimit",
		status: 403,
		enUS:   "Limit exceeded of udf number quota. quota is [%d] and already use [%d].",
		zhCN:   "函数个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientResourceLimit = &Error{
		code:   "QuotaInsufficientResourceLimit",
		status: 403,
		enUS:   "Limit exceeded of resource number quota. quota is [%d] and already use [%d].",
		zhCN:   "文件个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientResourceSize = &Error{
		code:   "QuotaInsufficientResourceSize",
		status: 403,
		enUS:   "Limit exceeded of single file size quota. quota is [%d] and already use [%d].",
		zhCN:   "单个文件大小超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientResourceSizeTotal = &Error{
		code:   "QuotaInsufficientResourceSizeTotal",
		status: 403,
		enUS:   "Limit exceeded of total file size quota. quota is [%d] and already use [%d].",
		zhCN:   "用户所有文件总大小超出配额限制. 配额 [%d], 已使用 [%d].",
	}

	QuotaInsufficientFlinkClusterLimit = &Error{
		code:   "QuotaInsufficientFlinkClusterLimit",
		status: 403,
		enUS:   "Limit exceeded of flink cluster number quota. quota is [%d] and already use [%d].",
		zhCN:   "Flink 计算集群个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}

	QuotaInsufficientFlinkClusterCU = &Error{
		code:   "QuotaInsufficientFlinkClusterCU",
		status: 403,
		enUS:   "Limit exceeded of single flink cu quota. quota is [%f] and already use [%f].",
		zhCN:   "单个 Flink 计算集群的 CU 个数超出配额限制. 配额 [%f], 已使用 [%f].",
	}

	QuotaInsufficientFlinkClusterCUTotal = &Error{
		code:   "QuotaInsufficientFlinkClusterCUTotal",
		status: 403,
		enUS:   "Limit exceeded of total cu quota. quota is [%f] and already use [%f].",
		zhCN:   "所有 Flink 计算集群的 CU 个数超出配额限制. 配额 [%f], 已使用 [%f].",
	}

	QuotaInsufficientNetworkLimit = &Error{
		code:   "QuotaInsufficientNetworkLimit",
		status: 403,
		enUS:   "Limit exceeded of network number quota. quota is [%d] and already use [%d].",
		zhCN:   "网络配置个数超出配额限制. 配额 [%d], 已使用 [%d].",
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
		status: http.StatusNoContent,
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

	EngineIncorrect = &Error{
		code:   "EngineIncorrect",
		status: 400,
		enUS:   "The engine[%s] not exist or status not in[%s].",
		zhCN:   "计算引擎[%s]不存在或者状态不是[%s].",
	}

	EngineNameInUse = &Error{
		code:   "EngineNameInUse",
		status: 400,
		enUS:   "The engine name[%s] is in use.",
		zhCN:   "计算集群名称[%s]已存在.",
	}

	NetworkNoAvailableIpRange = &Error{
		code:   "NoAvailableIpRange",
		status: 400,
		enUS:   "There is no available ip-range in vxnet[%s].",
		zhCN:   "网络[%s]里没有可用的IP段.",
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

// sign error
var (
	AccessKeyNotExists = &Error{
		code:   "AccessKeyNotExist",
		status: 403,
		enUS:   "The access key[%s] not exist.",
		zhCN:   "Access key[%s]不存在。",
	}
	UserNotExists = &Error{
		code:   "UserNotExist",
		status: 403,
		enUS:   "The user[%s] not exist.",
		zhCN:   "User[%s]不存在。",
	}
	ValidateSignatureFailed = &Error{
		code:   "ValidateSignatureFailed",
		status: 400,
		enUS:   "Validate signature failed, the signature[%s] not match [%s].",
		zhCN:   "签名验证失败, 签名[%s]与[%s]不匹配。",
	}

	MissingDateHeader = &Error{
		code:   "MissingDateHeader",
		status: 400,
		enUS:   "You must provide the Date or X-Date HTTP header.",
		zhCN:   "",
	}
	InvalidDateHeader = &Error{
		code:   "InvalidDateHeader",
		status: 400,
		enUS:   "The HTTP header Date or X-Date has wrong format.",
		zhCN:   "",
	}
	ExpiredSignature = &Error{
		code:   "ExpiredSignature",
		status: 401,
		enUS:   "The signature has been expired",
		zhCN:   "",
	}
	MissingAuthorizationHeader = &Error{
		code:   "MissingAuthorizationHeader",
		status: 400,
		enUS:   "You must provide the Authorization HTTP header.",
		zhCN:   "",
	}
	InvalidAuthorizationHeader = &Error{
		code:   "InvalidAuthorizationHeader",
		status: 400,
		enUS:   "The HTTP header Authorization has wrong format.",
		zhCN:   "",
	}
	UnsupportedSignatureVersion = &Error{
		code:   "UnsupportedSignVersion",
		status: 400,
		enUS:   "The signature version you used is not supported.",
		zhCN:   "",
	}
	SignatureNotMatch = &Error{
		code:   "SignatureNotMatch",
		status: 401,
		enUS:   "The request signature server calculated does not match the signature you provided.",
		zhCN:   "",
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

// resourcemanager error
var (
	FileSizeLimitExceededException = &Error{
		code:   "FileSizeLimitExceededException",
		status: 407,
		enUS:   "file size is large than 512mb",
		zhCN:   "文件大小超过 512mb",
	}
	InvalidFileSize = &Error{
		code:   "InvalidFileSize",
		status: 408,
		enUS:   "invalid file size [%s]",
		zhCN:   "文件大小不符",
	}
	HadoopClientCreateFailed = &Error{
		code:   "HadoopClientCreateFailed",
		status: 500,
		enUS:   "fail to init hadoop client",
		zhCN:   "hadoop客户端创建失败",
	}
)

// jobmanager error
var (
	CancelWithSavepointFailed = &Error{
		code:   "CancelWithSavepointFailed",
		status: 500,
		enUS:   "fail to cancel with savepoint reason [%s]",
		zhCN:   "取消任务失败: [%s]",
	}
	StreamJobSyntaxFailed = &Error{
		code:   "StreamJobSyntaxFailed",
		status: 400,
		enUS:   "%s",
		zhCN:   "%s",
	}
	ParseEngineFlinkUrlFailed = &Error{
		code:   "ParseEngineFlinkUrlFailed",
		status: 500,
		enUS:   "%s",
		zhCN:   "%s",
	}
)

// zeppelin error
var (
	InvalidateZeppelinUser = &Error{
		code:   "InvalidateZeppelinUser",
		status: 302,
		enUS:   "invalidate zeppelin user",
		zhCN:   "请重新登陆zeppelin",
	}

	CallZeppelinRestApiFailed = &Error{
		code:   "UnableCallZeppelinRestApi",
		status: 500,
		enUS:   "call zeppelin rest api failed, status: %s, statusText %s, message: %s",
		zhCN:   "zeppelin 服务异常",
	}

	ZeppelinNoteAlreadyExists = &Error{
		code:   "ZeppelinNoteAlreadyExists",
		status: 401,
		enUS:   "zeppelin note already exists",
		zhCN:   "notebook 已存在",
	}

	ZeppelinReturnStatusError = &Error{
		code:   "ZeppelinReturnStatusError",
		status: 500,
		enUS:   "submit zeppelin failed ,message: %s",
		zhCN:   "zeppelin 调用异常",
	}

	ZeppelinRunParagraphTimeout = &Error{
		code:   "ZeppelinRunParagraphTimeout",
		status: 500,
		enUS:   "the paragraph is not finished in %s seconds",
		zhCN:   "zeppelin paragraph 运行超时",
	}

	ZeppelinConfigureFailed = &Error{
		code:   "ZeppelinConfigureFailed",
		status: 500,
		enUS:   "fail to configure zeppelin session",
		zhCN:   "设置zeppelin config参数失败",
	}

	ZeppelinInitFailed = &Error{
		code:   "ZeppelinInitFailed",
		status: 500,
		enUS:   "fail to init zeppelin session %s",
		zhCN:   "zeppelin session 初始化失败",
	}

	RegisterUDFFailed = &Error{
		code:   "RegisterUDFFailed",
		status: 500,
		enUS:   "register udf failed",
		zhCN:   "注册udf失败",
	}

	ZeppelinSessionNotRunning = &Error{
		code:   "ZeppelinSessionNotRunning",
		status: 500,
		enUS:   "reconnect failed zeppelin session not running",
		zhCN:   "重连session失败",
	}

	ZeppelinGetJobIdFailed = &Error{
		code:   "ZeppelinGetJobIdFailed",
		status: 500,
		enUS:   "failed to get flink job id",
		zhCN:   "获取flink任务id失败",
	}

	ZeppelinParagraphRunError = &Error{
		code:   "ZeppelinRuntimeError",
		status: 500,
		enUS:   "job failed, reason: %s",
		zhCN:   "zeppelin任务运行失败",
	}
)

// flink rest error
var (
	FlinkRestError = &Error{
		code:   "FlinkRestError",
		status: 500,
		enUS:   "failed to ask request api,status: %s, statusText: %s, message:%s",
		zhCN:   "请求Flink APi失败",
	}

	FlinkJobNotExists = &Error{
		code:   "FlinkJobNotExists",
		status: 500,
		enUS:   "flink job not exists,jobName is %s",
		zhCN:   "Flink Job 不存在",
	}
)
