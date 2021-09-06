package constants

import "github.com/DataWorkbench/gproto/pkg/model"

type EngineRequestOptions struct {
	JobID       string `json:"jobid"`
	EngineID    string `json:"engineid"`
	WorkspaceID string `json:"workspaceid"`
	Parallelism int32  `json:"parallelism"`
	JobCU       int32  `json:"jobcu"`
	TaskCU      int32  `json:"taskcu"`
	TaskNum     int32  `json:"tasknum"`
	AccessKey   string `json:"accesskey"`
	SecretKey   string `json:"secretkey"`
	EndPoint    string `json:"endpoint"`
	//TODOHbaseHosts  []HostType `json:"hbasehosts"`
}

type EngineResponseOptions struct {
	EngineType      string `json:"enginetype"`
	EngineHost      string `json:"enginehost"`
	EnginePort      string `json:"engineport"`
	EngineExtension string `json:"engineextension"`
}

const (
	StatusFailed        = model.StreamFlowInst_Failed
	StatusFailedString  = "failed"
	StatusFinish        = model.StreamFlowInst_Succeed
	StatusFinishString  = "finish"
	StatusRunning       = model.StreamFlowInst_Running
	StatusRunningString = "running"
	JobSuccess          = "success"
	JobAbort            = "job abort"
	JobRunning          = "job running"
	jobError            = "error happend"
)
