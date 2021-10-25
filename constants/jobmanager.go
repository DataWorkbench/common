package constants

import "github.com/DataWorkbench/gproto/pkg/model"

const (
	StatusFailed           = model.StreamJobInst_Failed
	StatusFailedString     = "failed"
	StatusFinish           = model.StreamJobInst_Succeed
	StatusFinishString     = "finish"
	StatusRunning          = model.StreamJobInst_Running
	StatusRunningString    = "running"
	StatusTerminated       = model.StreamJobInst_Terminated
	StatusTerminatedString = "terminated"

	MessageUnknowState = "unknow state"
	MessageFinish      = "finish and success"
	MessageRunning     = "job running"
	MessageFailed      = "error happend"
	MessageTerminated  = "user terminated"

	JobCommandRun     = "run"
	JobCommandSyntax  = "syntax"
	JobCommandPreview = "preview"
	JobCommandExplain = "explain"
)
