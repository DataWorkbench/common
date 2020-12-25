package constants

// Node statue.
const (
	NodeStatusEnabled  int8 = iota + 1 // => "enabled"
	NodeStatusDisabled                 // => "disabled"
)

// Strategy of node task execute failure in a workflow.
const (
	NodeFailureStrategyNone   int8 = iota + 1 // => "none"
	NodeFailureStrategyIgnore                 // => "ignore"
)

// Defines the supported node type.
const (
	NodeTypeVirtual int8 = iota + 1 // => "virtual"
	NodeTypeShell                   // => "shell"
	NodeTypeFlinkJob
	NodeTypeFlinkSQL
)
