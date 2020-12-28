package constants

const (
	EngineTypeFlink      = "Flink"
	SourceTypeMysql      = "MySQL"
	SourceTypePostgreSQL = "PostgreSQL"
	SourceTypeKafka      = "Kafka"
	TableTypeDimension   = "d"
	TableTypeComment     = "c"
	CreatorWorkBench     = "workbench" //can't drop by custom,  workbench is auto create when spark/other engine created
	CreatorCustom        = "custom"
)
