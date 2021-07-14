package constants

const (
	ObsStreamFlowCycleInstTableName = "stream_workflow_cycle_instance"
	ObsDispatchedTaskCountTableName = "dispatched_task_count"
)

var DispatchedColumns = []string{
	"updated",
	"flow_count",
	"instance_count",
	"space_id",
}
