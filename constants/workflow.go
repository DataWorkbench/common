package constants

// Strategy of node task execute failure in a workflow.
const (
	ScheduleFailureStrategyContinue    int32 = iota + 1 // => "continue"
	ScheduleFailureStrategyTermination                  // => "termination"
)

// Strategy of schedule depends of workflow.
const (
	ScheduleConcurrencyPolicyAllow   int32 = iota + 1 // => "allow"
	ScheduleConcurrencyPolicyForbid                   // => "forbid"
	ScheduleConcurrencyPolicyReplace                  // => "replace"
)

// Strategy of schedule retry of workflow.
const (
	ScheduleRetryPolicyNone int32 = iota + 1
	ScheduleRetryPolicyAuto
)

// The environmental parameters for Flink.
type StreamFlowEnv struct {
	EngineId    string            `json:"engine_id"`
	Parallelism int32             `json:"parallelism"`
	JobCU       int32             `json:"job_cu"`
	TaskCU      int32             `json:"task_cu"`
	TaskNum     int32             `json:"task_num"`
	Custom      map[string]string `json:"custom"`
}

type StreamFlowSchedule struct {
	// Timestamp of start time of the validity period, unit in seconds.
	Started int64 `json:"started"`
	// Timestamp of end time of the validity period, unit in seconds.
	Ended int64 `json:"ended"`
	// Concurrency policy. 1 => "allow", 2 => "forbid", 3 => "replace"
	// - allow: Multiple task instances are allowed at the same time.
	// - forbid: No new instances will be created, and this schedule cycle will be skipped,
	// - replace: Force stop the old running instances and create new.
	ConcurrencyPolicy int32 `json:"concurrency_policy"`
	// Retry policy when task failed. 1 => "not retry" 2 => "auto retry"
	RetryPolicy   int32 `json:"retry_policy"`
	RetryLimit    int32 `json:"retry_limit"`
	RetryInterval int32 `json:"retry_interval"`
	Timeout       int32 `json:"timeout"`
	// Express is the standard unix crontab express.
	Express string `json:"express"`
}
