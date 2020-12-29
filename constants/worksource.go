package constants

const (
	EngineTypeFlink      = "Flink"
	SourceTypeMysql      = "MySQL"
	SourceTypePostgreSQL = "PostgreSQL"
	SourceTypeKafka      = "Kafka"
	TableTypeDimension   = "d"
	TableTypeCommon      = "c"
	CreatorWorkBench     = "workbench" //can't drop by custom,  workbench is auto create when spark/other engine created
	CreatorCustom        = "custom"
)

type SourceMysqlParams struct {
	User             string   `json:"user"`
	Password         string   `json:"password"`
	Host             string   `json:"host"`
	Port             int32    `json:"port"`
	Database         string   `json:"database"`
	ConnectorOptions []string `json:"connector_options"`
}

type SourcePostgreSQLParams struct {
	User             string   `json:"user"`
	Password         string   `json:"password"`
	Host             string   `json:"host"`
	Port             int32    `json:"port"`
	Database         string   `json:"database"`
	ConnectorOptions []string `json:"connector_options"`
}

type SourceKafkaParams struct {
	Host             string   `json:"host"`
	Port             int32    `json:"port"`
	ConnectorOptions []string `json:"connector_options"`
}

type FlinkTableDefine struct {
	SqlColumn        []string `json:"sql_column"`
	ConnectorOptions []string `json:"connector_options"`
}
