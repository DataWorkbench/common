package constants

const (
	InstanceTriggerSchedule      int32 = iota + 1 // => "schedule"
	InstanceTriggerManual                         // => "manual" && "test"
	InstanceTriggerSupplementary                  // => "supplementary data"
)

const (
	InstanceStatePending   int32 = iota + 1 // => "pending"
	InstanceStateRunning                    // => "running"
	InstanceStateSuspended                  // => "suspended"
	InstanceStateStopped                    // => "stopped"
	InstanceStateSucceed                    // => "succeed" -- jobmanager
	InstanceStateFailed                     // => "failed"  -- jobmanager
)
