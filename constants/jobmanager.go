package constants

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
}

type EngineResponseOptions struct {
	EngineType      string `json:"enginetype"`
	EngineHost      string `json:"enginehost"`
	EnginePort      string `json:"engineport"`
	EngineExtension string `json:"engineextension"`
}

const (
	StatusFailed        = InstanceStateFailed
	StatusFailedString  = "failed"
	StatusFinish        = InstanceStateSucceed
	StatusFinishString  = "finish"
	StatusRunning       = InstanceStateRunning
	StatusRunningString = "running"
	JobSuccess          = "success"
	JobAbort            = "job abort"
	JobRunning          = "job running"
	jobError            = "error happend"
)
