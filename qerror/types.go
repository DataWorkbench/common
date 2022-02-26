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

	// ResourceNotActive render message if accessing resource is inactive and operation not allows.
	ResourceNotActive = &Error{
		code:   "ResourceNotActive",
		status: 403,
		enUS:   "The resource [%s] you are accessing does not active.",
		zhCN:   "资源 [%s] 已被禁用, 禁止访问.",
	}

	// ResourceNotExists render message if accessing resource not exists.
	ResourceNotExists = &Error{
		code:   "ResourceNotExists",
		status: 404,
		enUS:   "The resource [%s] you are accessing does not exist.",
		zhCN:   "资源 [%s] 不存在",
	}

	// ResourceAlreadyExists render message if creating resource already exists.
	ResourceAlreadyExists = &Error{
		code:   "ResourceAlreadyExists",
		status: 409,
		enUS:   "The resource [%s] you are using already exists.",
		zhCN:   "资源 [%s] 已存在",
	}

	// ResourceIsInUsing render message if be deletion resource is using by other module.
	ResourceIsInUsing = &Error{
		code:   "ResourceIsInUsing",
		status: 403,
		enUS:   "The resource [%s] is used by [%s]",
		zhCN:   "资源 [%s] 正在被 [%s] 使用中",
	}
)

// parameters error
var (
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
		enUS:   "Found unknown field [%s] in you request",
		zhCN:   "请求中包含不支持的字段[%s]",
	}
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
	ParameterValidationError = &Error{
		code:   "ParameterValidationError",
		status: 400,
		enUS:   "%s",
		zhCN:   "%s",
	}
	InvalidRequest = &Error{
		code:   "InvalidRequest",
		status: 400,
		enUS:   "%s",
		zhCN:   "%s",
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
		enUS:   "You must provide the Date or X-Date in HTTP header.",
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

// general error.
var (
	InvalidName = &Error{
		code:   "InvalidName",
		status: 400,
		enUS:   "The name is invalid, accepts 0~9、a~z、 A~Z、_ and can't begin or end with _.",
		zhCN:   "名称不符合要求, 只允许数字,大小写字母和下划线, 并且不能以下划线开头或者结尾.",
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
	QuotaInsufficientFileLimit = &Error{
		code:   "QuotaInsufficientFileLimit",
		status: 403,
		enUS:   "Limit exceeded of resource number quota. quota is [%d] and already use [%d].",
		zhCN:   "文件个数超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientFileSize = &Error{
		code:   "QuotaInsufficientFileSize",
		status: 403,
		enUS:   "Limit exceeded of single file size quota. quota is [%d] and already use [%d].",
		zhCN:   "单个文件大小超出配额限制. 配额 [%d], 已使用 [%d].",
	}
	QuotaInsufficientFileSizeTotal = &Error{
		code:   "QuotaInsufficientFileSizeTotal",
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
		enUS:   "Limit exceeded of single flink cu quota. quota is [%0.1f] and already use [%0.1f].",
		zhCN:   "单个 Flink 计算集群的 CU 个数超出配额限制. 配额 [%0.1f], 已使用 [%0.1f].",
	}

	QuotaInsufficientFlinkClusterCUTotal = &Error{
		code:   "QuotaInsufficientFlinkClusterCUTotal",
		status: 403,
		enUS:   "Limit exceeded of total cu quota. quota is [%0.1f] and already use [%0.1f].",
		zhCN:   "所有 Flink 计算集群的 CU 个数超出配额限制. 配额 [%0.1f], 已使用 [%0.1f].",
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
//RegionNotSpecified = &Error{
//	code:   "RegionNotSpecified",
//	status: 400,
//	enUS:   "A valid region id must be specified in you request path",
//	zhCN:   "无效的请求, 未指定 regionId",
//}
//RegionNotExists = &Error{
//	code:   "RegionNotExists",
//	status: 404,
//	enUS:   "The region [%s] you access not exists.",
//	zhCN:   "访问的区域[%s]不存在.",
//}
//RegionAccessDenied = &Error{
//	code:   "RegionAccessDenied",
//	status: 403,
//	enUS:   "The user [%s] is not allowed to access region [%s].",
//	zhCN:   "用户 [%s] 没有访问区域 [%s] 的权限",
//}
)

// workspace error
var (
	SpaceProhibitDelete = &Error{
		code:   "SpaceProhibitDelete",
		status: 403,
		enUS:   "The workspace [%s] cannot be deleted, Please delete all flink clusters in the space first.",
		zhCN:   "工作空间 [%s] 不能被删除, 请先删除空间内的所有计算集群.",
	}

//SpaceNotExists = &Error{
//	code:   "SpaceNotExists",
//	status: 404,
//	enUS:   "The workspace [%s] does not exists.",
//	zhCN:   "工作空间[%s]不存在.",
//}
//SpaceAlreadyExists = &Error{
//	code:   "SpaceAlreadyExists",
//	status: 409,
//	enUS:   "The workspace name [%s] has been used.",
//	zhCN:   "工作空间名称[%s]已被使用.",
//}
//SpaceNotActive = &Error{
//	code:   "SpaceNotActive",
//	status: 403,
//	enUS:   "The workspace [%s] does not active.",
//	zhCN:   "工作空间[%s]已被禁用",
//}
)

// member error
var (
//MemberNotExists = &Error{
//	code:   "MemberNotExists",
//	status: 404,
//	enUS:   "The member [%s] does not exists.",
//	zhCN:   "成员[%s]不存在.",
//}
//MemberAlreadyExists = &Error{
//	code:   "MemberAlreadyExists",
//	status: 409,
//	enUS:   "The member [%s] has been exists.",
//	zhCN:   "空间成员[%s]已存在.",
//}
//SpaceOwnerCannotBeDeletion = &Error{
//	code:   "SpaceOwnerCannotDeletion",
//	status: 403,
//	enUS:   "The member of workspace owner cannot be deleted.",
//	zhCN:   "空间所有者不允许被删除.",
//}
//SpaceOwnerCannotBeUpdated = &Error{
//	code:   "SpaceOwnerCannotBeUpdated",
//	status: 403,
//	enUS:   "The member of workspace owner cannot be updated.",
//	zhCN:   "空间所有者不允许被修改",
//}
)

// stream job error.
var (
	StreamJobScheduleNotSet = &Error{
		code:   "ScheduleNotSet",
		status: 400,
		enUS:   "The stream job [%s] not set schedule properties.",
		zhCN:   "实时作业[%s]未设置调度属性",
	}
	StreamJobCodeNotSet = &Error{
		code:   "StreamJobCodeNotSet",
		status: 400,
		enUS:   "The stream job [%s] not set node task.",
		zhCN:   "实时作业[%s]未设置节点任务",
	}
	StreamJobArgsNotSet = &Error{
		code:   "ArgsNotSet",
		status: 400,
		enUS:   "The stream job [%s] not set environmental parameters ",
		zhCN:   "实时作业[%s]未设置环境参数",
	}
)

// sync job error.
var (
	SyncJobScheduleNotSet = &Error{
		code:   "ScheduleNotSet",
		status: 400,
		enUS:   "The sync job [%s] not set schedule properties.",
		zhCN:   "同步作业[%s]未设置调度属性",
	}
	SyncJobArgsNotSet = &Error{
		code:   "ArgsNotSet",
		status: 400,
		enUS:   "The sync job [%s] not set config parameters ",
		zhCN:   "同步作业[%s]未设置配置参数",
	}
)

// instance error
var (
//InstanceNotExists = &Error{
//	code:   "InstanceNotExists",
//	status: 404,
//	enUS:   "The instance [%s] does not exists.",
//	zhCN:   "任务实例[%s]不存在.",
//}
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
	//InvalidSourceName = &Error{
	//	code:   "InvalidSourceName",
	//	status: http.StatusInternalServerError,
	//	enUS:   "invalid name. can't use '.'",
	//	zhCN:   "无效名字，不能使用'.'",
	//}
	//SourceIsDisable = &Error{
	//	code:   "SourceIsDisable",
	//	status: http.StatusInternalServerError,
	//	enUS:   "source is disable",
	//	zhCN:   "数据源是禁用状态",
	//}
	//ConnectSourceFailed = &Error{
	//	code:   "ConnectSourceFailed",
	//	status: http.StatusNoContent,
	//	enUS:   "this source can't connect[%d]",
	//	zhCN:   "数据源无法连接[%d]",
	//}
)

// udfmanager
var (
//InvalidUDFName = &Error{
//	code:   "InvalidUDFName",
//	status: http.StatusInternalServerError,
//	enUS:   "invalid name. can't use '.'",
//	zhCN:   "无效名字，不能使用'.'",
//}
//NotSupportUDFType = &Error{
//	code:   "NotSupportUDFType",
//	status: http.StatusInternalServerError,
//	enUS:   "not support UDF type[%d]",
//	zhCN:   "不支持的自定义函数类型[%d]",
//}
)

// enginemanager error
var (
	//EngineNotExist = &Error{
	//	code:   "EngineNotExist",
	//	status: 404,
	//	enUS:   "The engine[%s] not exist.",
	//	zhCN:   "计算引擎[%s]不存在.",
	//}
	//EngineIncorrect = &Error{
	//	code:   "EngineIncorrect",
	//	status: 400,
	//	enUS:   "The engine[%s] not exist or status not in[%s].",
	//	zhCN:   "计算引擎[%s]不存在或者状态不是[%s].",
	//}
	//EngineNameInUse = &Error{
	//	code:   "EngineNameInUse",
	//	status: 400,
	//	enUS:   "The engine name[%s] is in use.",
	//	zhCN:   "计算集群名称[%s]已存在.",
	//}

	NetworkNoAvailableIpRange = &Error{
		code:   "NoAvailableIpRange",
		status: 400,
		enUS:   "There is no available ip-range in vxnet[%s].",
		zhCN:   "网络[%s]里没有可用的IP段.",
	}
)

// notifier error
var (
//SendingNotificationFailed = &Error{
//	code:   "SendingNotificationFailed",
//	status: 500,
//	enUS:   "failed to send notifier post to user [%s]",
//	zhCN:   "向用户[%s]发送通知失败.",
//}
//
//NotificationChannelConfigInvalid = &Error{
//	code:   "NotificationChannelConfigInvalid",
//	status: 500,
//	enUS:   "notifier channel config not valid",
//	zhCN:   "通知渠道配置不合法",
//}
)

// logmanager error
var (
//LogFileNotExists = &Error{
//	code:   "LogFileNotExists",
//	status: 405,
//	enUS:   "Log file [%s] not exists",
//	zhCN:   "日志文件 [%s] 不存在",
//}
//
//RequestForFlinkFailed = &Error{
//	code:   "RequestForFlinkFailed",
//	status: 406,
//	enUS:   "fail to request flink api [%s]",
//	zhCN:   "调用Flink Web Api [%s] 失败",
//}
)

// resourcemanager error
var (
	//FileSizeLimitExceededException = &Error{
	//	code:   "FileSizeLimitExceededException",
	//	status: 407,
	//	enUS:   "file size is large than 512mb",
	//	zhCN:   "文件大小超过 512mb",
	//}
	//InvalidFileSize = &Error{
	//	code:   "InvalidFileSize",
	//	status: 408,
	//	enUS:   "invalid file size [%s]",
	//	zhCN:   "文件大小不符",
	//}

	// FIXME: remove it.
	HadoopClientCreateFailed = &Error{
		code:   "HadoopClientCreateFailed",
		status: 500,
		enUS:   "fail to init hadoop client",
		zhCN:   "hadoop客户端创建失败",
	}
)

// jobmanager error
var (
	StreamJobSyntaxFailed = &Error{
		code:   "StreamJobSyntaxFailed",
		status: 400,
		enUS:   "%s",
		zhCN:   "%s",
	}
	//CancelWithSavepointFailed = &Error{
	//	code:   "CancelWithSavepointFailed",
	//	status: 500,
	//	enUS:   "fail to cancel with savepoint reason [%s]",
	//	zhCN:   "取消任务失败: [%s]",
	//}
	// FIXME: removed it.
	ParseEngineFlinkUrlFailed = &Error{
		code:   "ParseEngineFlinkUrlFailed",
		status: 500,
		enUS:   "%s",
		zhCN:   "%s",
	}
)

// zeppelin error
var (
	ZeppelinInitFailed = &Error{
		code:   "ZeppelinInitFailed",
		status: 500,
		enUS:   "fail to init zeppelin session %s",
		zhCN:   "zeppelin session 初始化失败",
	}

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
	//FlinkJobNotExists = &Error{
	//	code:   "FlinkJobNotExists",
	//	status: 500,
	//	enUS:   "flink job not exists,jobName is %s",
	//	zhCN:   "Flink Job 不存在",
	//}
)
