package constants

import "github.com/DataWorkbench/common/utils"

const (
	EngineFlinkDefaultId = IdPrefixEngine + "0000000000000000"
)

const (
	// ret code
	RetSucc = 0

	// Table const
	TableEngine        = "engine"
	TableEngineInBuild = "engine_in_build_info"

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

	// engine helm chart
	HelmFlink = "flink-1.12.3.tgz"

)

var InbuildEngineTypes = utils.StrArray{
	EngineTypeFlink,
}

var EngineTypeHelmChartMap = map[string]string{
	EngineTypeFlink: HelmFlink,
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
