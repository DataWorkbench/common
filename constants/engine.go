package constants

const (
	EngineFlinkDefaultId = IdPrefixEngine + "0000000000000000"
)

const (

	// Engine status
	EngineStatusInit    = "initialized"
	EngineStatusFail    = "failed"
	EngineStatusActive  = "active"
	EngineStatusDeleted = "deleted"

	EngineStatusDisable = "disable"
	EngineStatusEnable  = "enable"

	EngineTransactionStatusCreating = "creating"
	EngineTransactionStatusUpdating = "updating"
	EngineTransactionStatusDeleting = "deleting"

	// Engine server Type
	EngineTypeFlink = "flink"
	EngineTypeSpark = "spark"
)
