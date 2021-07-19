package constants

const (
	SourceTypeMysql        = "MySQL"
	SourceTypePostgreSQL   = "PostgreSQL"
	SourceTypeKafka        = "Kafka"
	SourceTypeS3           = "S3"
	SourceTypeClickHouse   = "ClickHouse"
	SourceTypeHbase        = "Hbase"
	DirectionSource        = "s"
	DirectionDestination   = "d"
	SourceConnectedSuccess = "t"
	SourceConnectedFailed  = "f"
	CreatorWorkBench       = "workbench" //can't drop by custom,  workbench is auto create when spark/other engine created
	CreatorCustom          = "custom"
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

type SourceS3Params struct {
	AccessKey string `json:"accesskey"`
	SecretKey string `json:"secretkey"`
	EndPoint  string `json:"endpoint"`
}

type SourceClickHouseParams struct {
	User             string   `json:"user"`
	Password         string   `json:"password"`
	Host             string   `json:"host"`
	Port             int32    `json:"port"`
	Database         string   `json:"database"`
	ConnectorOptions []string `json:"connector_options"`
}

type SourceHbaseParams struct {
	Zookeeper string `json:"zookeeper"`
	Znode     string `json:"znode"`
	Hosts     string `json:"hosts"`
}

type FlinkTableDefineHbase struct {
	SqlColumn        []string `json:"sql_column"`
	ConnectorOptions []string `json:"connector_options"`
}

type FlinkTableDefineClickHouse struct {
	SqlColumn        []string `json:"sql_column"`
	ConnectorOptions []string `json:"connector_options"`
}

type FlinkTableDefineS3 struct {
	SqlColumn        []string `json:"sql_column"`
	Path             string   `json:"path"`
	Format           string   `json:"format"`
	ConnectorOptions []string `json:"connector_options"`
}

type FlinkTableDefineKafka struct {
	SqlColumn        []string `json:"sql_column"`
	Topic            string   `json:"topic"`
	Format           string   `json:"format"`
	ConnectorOptions []string `json:"connector_options"`
}

type FlinkTableDefineMysql struct {
	SqlColumn        []string `json:"sql_column"`
	ConnectorOptions []string `json:"connector_options"`
}

type FlinkTableDefinePostgreSQL struct {
	SqlColumn        []string `json:"sql_column"`
	ConnectorOptions []string `json:"connector_options"`
}
