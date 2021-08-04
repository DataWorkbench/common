package constants

import "github.com/DataWorkbench/common/utils"

const (
	EngineFlinkDefaultId = IdPrefixEngine + "0000000000000000"
)

const (
	// ret code
	RetSucc = 0

	// Table const
	EngineTable        = "engine"
	EngingInBuildTable = "engine_in_build_info"

	// Engine status
	EngineStatusDisable = "disable"
	EngineStatusEnable = "enable"
	EngineStatusDeleted = "deleted"

	EngineTransitionStatusCreating = "creating"
	EngineTransitionStatusUpdating = "updating"
	EngineTransitionStatusDeleting = "deleting"

	// Engine server Type
	EngineTypeFlink = "flink"
	EngineTypeSpark = "spark"
)

var EngineTypes = utils.StrArray{
	EngineTypeFlink,
	EngineTypeSpark,
}

var EngineColumns = []string{
	"id",
	"name",
	"owner",
	"desc",
	"url",
	"server_type",
	"status",
	"is_inbuilt",
	"created",
	"updated",
}
