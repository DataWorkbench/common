package constants

const (
	InstanceStatePending    int32 = iota + 1 // => "pending"
	InstanceStateRunning                     // => "running" -- jobmanager. RunJob return.
	InstanceStateSuspended                   // => "suspended"
	InstanceStateTerminated                  // => "terminated"
	InstanceStateRetrying                    // => "retrying" - failed and retry
	InstanceStateTimeout                     // => "timeout" - task execute timeout.
	InstanceStateSucceed                     // => "succeed" -- jobmanager. GetJobState return
	InstanceStateFailed                      // => "failed"  -- jobmanager. GetJobState return
)
