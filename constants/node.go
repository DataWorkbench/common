package constants

// Node statue.
const (
	NodeStatusEnabled  int32 = iota + 1 // => "enabled"
	NodeStatusDisabled                  // => "disabled"
)

// Strategy of node task execute failure in a workflow.
const (
	NodeFailureStrategyNone   int32 = iota + 1 // => "none"
	NodeFailureStrategyIgnore                  // => "ignore"
)

// Defines the supported node type.
const (
	NodeTypeVirtual int32 = iota + 1 // => "virtual"
	NodeTypeShell                    // => "shell"
	NodeTypeFlinkJob
	NodeTypeFlinkSSQL
)

const (
	MainRunQuote = "$qc$"
)

// Defines of NodeTypeFlinkSSQL.
type FlinkSSQL struct {
	Tables      []string `json:"tables"`
	Funcs       []string `json:"funcs"`
	Parallelism int32    `json:"parallelism"`
	JobMem      int32    `json:"job_mem"` // in MB
	JobCpu      float32  `json:"job_cpu"`
	TaskCpu     float32  `json:"task_cpu"`
	TaskMem     int32    `json:"task_mem"` // in MB
	TaskNum     int32    `json:"task_num"`
	MainRun     string   `json:"main_run"` //AccessKey, SecretKey, EndPoint is in sourcemanager
}

// Defines of NodeTypeFlinkJob.
type FlinkJob struct {
	Parallelism int32   `json:"parallelism"`
	JobMem      int32   `json:"job_mem"` // in MB
	JobCpu      float32 `json:"job_cpu"`
	TaskCpu     float32 `json:"task_cpu"`
	TaskMem     int32   `json:"task_mem"` // in MB
	TaskNum     int32   `json:"task_num"`
	JarArgs     string  `json:"jar_args"`  // allow regex `^[a-zA-Z0-9_/. ]+$`
	JarEntry    string  `json:"jar_entry"` // allow regex `^[a-zA-Z0-9_/. ]+$`
	MainRun     string  `json:"main_run"`
	AccessKey   string  `json:"accesskey"`
	SecretKey   string  `json:"secretkey"`
	EndPoint    string  `json:"endpoint"`
}
