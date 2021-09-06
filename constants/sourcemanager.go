package constants

const (
	SourceTypeMysql      = "MySQL"
	SourceTypePostgreSQL = "PostgreSQL"
	SourceTypeKafka      = "Kafka"
	SourceTypeS3         = "S3"
	SourceTypeClickHouse = "ClickHouse"
	SourceTypeHbase      = "Hbase"
	SourceTypeFtp        = "Ftp"
	SourceTypeHDFS       = "HDFS"

	DirectionSource      = "source"
	DirectionDestination = "distination"
	DirectionDimension   = "dimension"

	SourceEnableState  = "enable"
	SourceDisableState = "disable"

	SourceConnectedSuccess = "successful"
	SourceConnectedFailed  = "failed"
)
