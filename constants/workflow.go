package constants

// Workflow Release Status.
const (
	FlowReleaseStatusActive    int32 = iota + 1 // => "active"
	FlowReleaseStatusSuspended                  // => "suspended"
)

// Workflow Type.
const (
	FlowTypeStreamSQL      int32 = iota + 1 // => "stream works with SQL"
	FlowTypeStreamJAR                       // => "stream works with JAR ball"
	FlowTypeStreamOperator                  // => "stream works with operator choreography".
)

//// Workflow priority.
//const (
//	FlowPriorityHighest int32 = iota + 1 // => "highest"
//	FlowPriorityHigh                     // => "high"
//	FlowPriorityMedium                   // => "medium"
//	FlowPriorityLow                      // => "low"
//	FlowPriorityLowest                   // => "lowest"
//)

//// Strategy of node task execute failure in a workflow.
//const (
//	ScheduleFailureStrategyContinue    int32 = iota + 1 // => "continue"
//	ScheduleFailureStrategyTermination                  // => "termination"
//)

// Strategy of schedule depends of workflow.
const (
	ScheduleDependStrategyNone int32 = iota + 1 // => "none"
	ScheduleDependStrategyLast                  // => "last"
	ScheduleDependStrategyStop                  // => "stop"
)

// Strategy of schedule retry of workflow.
const (
	ScheduleRetryStrategyNone int32 = iota + 1
	ScheduleRetryStrategyAuto
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

// The environmental parameters for Flink.
type StreamFlowEnv struct {
	EngineId    string            `json:"engine_id"`
	Parallelism int32             `json:"parallelism"`
	JobMem      int32             `json:"job_mem"` // in MB
	JobCpu      float32           `json:"job_cpu"`
	TaskCpu     float32           `json:"task_cpu"`
	TaskMem     int32             `json:"task_mem"` // in MB
	TaskNum     int32             `json:"task_num"`
	Custom      map[string]string `json:"custom"`
}

type StreamFlowSchedule struct {
	// Timestamp of start time of the validity period, unit in seconds.
	Started int64 `json:"started"`
	// Timestamp of end time of the validity period, unit in seconds.
	Ended int64 `json:"ended"`
	// Strategy of dependency. 1 => "none", 2 => "last".
	DependStrategy int32 `json:"depend_strategy"`
	RetryStrategy  int32 `json:"retry_strategy"`
	RetryLimit     int32 `json:"retry_limit"`
	RetryInterval  int32 `json:"retry_interval"`
	Timeout        int32 `json:"timeout"`
	// Express is the standard unix crontab express.
	Express string `json:"express"`
}
