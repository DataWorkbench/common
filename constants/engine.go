package constants

const (
	EngineFlinkDefaultId = IdPrefixEngine + "0000000000000000"
)

const (
	// ret code
	RetSucc = 0

	// Table const
	EngineIDPrefix = "eng-"
	EngineTable    = "engine"

	// Engine status
	EngineStatusDisable int8 = iota
	EngineStatusEnable
	EngineStatusDeleted

	// Engine server Type
	ServerTypeFlink = "flink"
	ServerTypeSpark = "spark"

	// Engine build Type
	InBuilt  = 1
	External = 0
)

var ServerTypes = []string{
	ServerTypeFlink,
	ServerTypeSpark,
}

var StatusCode = map[int8]string{
	EngineStatusEnable:  "enable",
	EngineStatusDisable: "disable",
	EngineStatusDeleted: "deleted",
}

var BuildTypeCode = map[int8]string{
	InBuilt:  "inbuilt",
	External: "external",
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
