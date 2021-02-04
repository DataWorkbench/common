package constants

type EngineRequestOptions struct {
	JobID       string  `json:"jobid"`
	WorkspaceID string  `json:"workspaceid"`
	Parallelism int32   `json:"parallelism"`
	JobMem      int32   `json:"job_mem"` // in MB
	JobCpu      float32 `json:"job_cpu"`
	TaskCpu     float32 `json:"task_cpu"`
	TaskMem     int32   `json:"task_mem"` // in MB
	TaskNum     int32   `json:"task_num"`
	AccessKey   string  `json:"accesskey"`
	SecretKey   string  `json:"secretkey"`
	EndPoint    string  `json:"endpoint"`
}

type EngineResponseOptions struct {
	EngineType      string `json:"enginetype"`
	EngineHost      string `json:"enginehost"`
	EnginePort      string `json:"engineport"`
	EngineExtension string `json:"engineextension"`
}
