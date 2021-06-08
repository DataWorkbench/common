package constants

const (
	InstanceTriggerSchedule      int32 = iota + 1 // => "cycle"
	InstanceTriggerManual                         // => "manual" && "test"
	InstanceTriggerSupplementary                  // => "supplementary data"
)

const (
	InstanceStatePending    int32 = iota + 1 // => "pending"
	InstanceStateRunning                     // => "running" -- jobmanager. RunJob return.
	InstanceStateSuspended                   // => "suspended"
	InstanceStateTerminated                  // => "terminated"
	InstanceStateSucceed                     // => "succeed" -- jobmanager. Report return
	InstanceStateFailed                      // => "failed"  -- jobmanager. Report return
)
