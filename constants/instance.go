package constants

const (
	InstanceStatePending    int32 = iota + 1 // => "pending"
	InstanceStateRunning                     // => "running" -- jobmanager. RunJob return.
	InstanceStateRetrying                    // => "retrying" -- failed and retry
	InstanceStateSuspended                   // => "suspended" -- unused now.
	InstanceStateTerminated                  // => "terminated"
	InstanceStateTimeout                     // => "timeout" - task execute timeout.
	InstanceStateSucceed                     // => "succeed" -- jobmanager. GetJobState return
	InstanceStateFailed                      // => "failed"  -- jobmanager. GetJobState return
)
