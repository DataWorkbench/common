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
	SyntaxCheckCommand = "syntaxcheck"
)

type FlinkNode struct {
	StreamSql bool       `json:"stream_sql"` //true is flink stream sql. false is flink batch sql
	Env       JSONString `json:"env"`
	Nodes     JSONString `json:"nodes"`
}
