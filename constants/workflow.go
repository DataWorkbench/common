package constants

// Workflow priority.
const (
	FlowPriorityHighest int32 = iota + 1 // => "highest"
	FlowPriorityHigh                     // => "high"
	FlowPriorityMedium                   // => "medium"
	FlowPriorityLow                      // => "low"
	FlowPriorityLowest                   // => "lowest"
)

// Strategy of node task execute failure in a workflow.
const (
	FlowFailureStrategyContinue int32 = iota + 1 // => "continue"
	FlowFailureStrategySuspend                   // => "suspend"
)

// Strategy of schedule depends of workflow.
const (
	FlowDependStrategyNone int32 = iota + 1 // => "none"
	FlowDependStrategyLast                  // => "last"
)

// Strategy of schedule.
const (
	FlowScheduleStrategyLoop int32 = iota + 1 // => "loop"
)

// Strategy of notify of workflow.
const (
	FlowNotifyStrategyFlowStarted int32 = iota + 1
	FlowNotifyStrategyFlowSucceed
	FlowNotifyStrategyFlowFailed
	FlowNotifyStrategyNodeStarted
	FlowNotifyStrategyNodeSucceed
	FlowNotifyStrategyNodeRetried
	FlowNotifyStrategyNodeFailed
)

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
	NodeTypeFlinkSQL
)
