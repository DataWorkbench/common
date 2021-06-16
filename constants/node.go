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
	NodeTypeFlink
)

const (
	RunCommand         = "run"
	PreviewCommand     = "preview"
	ExplainCommand     = "explain"
	SyntaxCheckCommand = "syntax"
)

type FlinkNode struct {
	Command     string     `json:"command"`
	StreamSql   bool       `json:"stream_sql"` //true is flink stream sql. false is flink batch sql
	Parallelism int32      `json:"parallelism"`
	JobMem      int32      `json:"job_mem"`  // use in serverless engine // MB
	JobCpu      float32    `json:"job_cpu"`  // use in serverless engine
	TaskCpu     float32    `json:"task_cpu"` // use in serverless engine
	TaskMem     int32      `json:"task_mem"` // use in serverless engine // MB
	TaskNum     int32      `json:"task_num"` // use in serverless engine
	Nodes       JSONString `json:"nodes"`
}
