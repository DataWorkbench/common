package constants

// Workflow priority.
const (
	FlowPriorityHighest int8 = iota + 1 // => "highest"
	FlowPriorityHigh                    // => "high"
	FlowPriorityMedium                  // => "medium"
	FlowPriorityLow                     // => "low"
	FlowPriorityLowest                  // => "lowest"
)

// Strategy of node task execute failure in a workflow.
const (
	FlowFailureStrategyContinue int8 = iota + 1 // => "continue"
	FlowFailureStrategyClosure                  // => "closure"
)

// Strategy of schedule depends of workflow.
const (
	FlowDependStrategyNone int8 = iota + 1 // => "none"
	FlowDependStrategyLast                 // => "last"
)

// Strategy of schedule.
const (
	FlowScheduleStrategyLoop int8 = iota + 1 // => "loop"
	FlowScheduleStrategyOnce                 // => "once"
)

// Strategy of notify of workflow.
const (
	FlowNotifyStrategyFlowStarted int8 = iota + 1
	FlowNotifyStrategyFlowSucceed
	FlowNotifyStrategyFlowFailed
	FlowNotifyStrategyNodeStarted
	FlowNotifyStrategyNodeSucceed
	FlowNotifyStrategyNodeRetried
	FlowNotifyStrategyNodeFailed
)
