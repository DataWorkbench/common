package constants

const (
	FlinkVersion_011203_0211 = "flink-1.12.3-scala_2.11"
)

const (
	FlinkConnectorMySQL          = "flink-connector-mysql"
	FlinkConnectorMySQLCDC       = "flink-connector-mysql-cdc"
	FlinkConnectorPostgresSQL    = "flink-connector-postgresql"
	FlinkConnectorKafka          = "flink-connector-kafka"
	FlinkConnectorHbase          = "flink-connector-hbase"
	FlinkConnectorClickhouse     = "flink-connector-clickhouse"
	FlinkConnectorElasticsearch7 = "flink-connector-elasticsearch7"
)

// FlinkConnectorLists represents the list of built-in connectors.
var FlinkConnectorLists = []string{
	FlinkConnectorMySQL,
	FlinkConnectorMySQLCDC,
	FlinkConnectorPostgresSQL,
	FlinkConnectorKafka,
	FlinkConnectorHbase,
	FlinkConnectorClickhouse,
	FlinkConnectorElasticsearch7,
}

// FlinkConnectorJarMap represents the map to connectors and jar ball.
// {"FlinkVersion": {"ConnectorName": {"Jar1", "Jar2"}} }
var FlinkConnectorJarMap = map[string]map[string][]string{
	FlinkVersion_011203_0211: {
		// /flinkc/buildin/connectors/mysql.jar
		FlinkConnectorMySQL:          {"mysql-connector-java-8.0.21.jar", "flink-connector-jdbc_2.11-1.12.3.jar"},
		FlinkConnectorMySQLCDC:       {"flink-connector-mysql-cdc-1.3.0.jar", "flink-connector-jdbc_2.11-1.12.3.jar"},
		FlinkConnectorPostgresSQL:    {"postgresql-42.2.18.jar", "flink-connector-jdbc_2.11-1.12.3.jar"},
		FlinkConnectorKafka:          {"flink-sql-connector-kafka_2.11-1.12.3.jar"},
		FlinkConnectorHbase:          {"flink-sql-connector-hbase-2.2_2.11-1.12.3.jar"},
		FlinkConnectorClickhouse:     {"flink-connector-clickhouse-1.0.0.jar", "flink-connector-jdbc_2.11-1.12.3.jar"},
		FlinkConnectorElasticsearch7: {"flink-sql-connector-elasticsearch7_2.11-1.12.3.jar"},
	},
}