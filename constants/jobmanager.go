package constants

import "github.com/DataWorkbench/gproto/pkg/model"

const (
	StatusFailed           = model.StreamFlowInst_Failed
	StatusFailedString     = "failed"
	StatusFinish           = model.StreamFlowInst_Succeed
	StatusFinishString     = "finish"
	StatusRunning          = model.StreamFlowInst_Running
	StatusRunningString    = "running"
	StatusTerminated       = model.StreamFlowInst_Terminated
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
