package dsutils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel/pbdatasource"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"github.com/Shopify/sarama"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"gopkg.in/mgo.v2"
	"time"

	"github.com/dazheng/gohive"
	elastic6 "github.com/olivere/elastic/v6"
	elastic7 "github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DescribeDatasourceTablesMySQL get a table list of type MySQL.
func DescribeDatasourceTablesMySQL(ctx context.Context, url *pbdatasource.MySQLURL) (items []string, err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		url.User, url.Password, url.Host, url.Port, url.Database,
	)

	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	defer func() {
		// close the connections.
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}()

	err = db.Raw("select table_name as item from information_schema.tables where  table_schema = ?",
		url.Database).Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return
}

// DescribeDatasourceTablesPostgresSQL get a table list of type PostgresSQL.
func DescribeDatasourceTablesPostgresSQL(ctx context.Context, url *pbdatasource.PostgreSQLURL) (items []string,
	err error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s ",
		url.User, url.Password, url.Host, url.Port, url.Database,
	)
	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	defer func() {
		// close the connections.
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}()

	err = db.Raw("SELECT n.nspname || '.' || c.relname as item FROM pg_catalog.pg_class c LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace WHERE c.relkind IN ('r','p') AND n.nspname <> 'pg_catalog' AND n.nspname <> 'information_schema' AND n.nspname !~ '^pg_toast' AND pg_catalog.pg_table_is_visible(c.oid);").
		Scan(&items).Error
	if err != nil {
		return nil, err
	}
	return
}

// DescribeDatasourceTablesClickHouse get a table list of type ClickHouse.
func DescribeDatasourceTablesClickHouse(ctx context.Context, url *pbdatasource.ClickHouseURL) (items []string,
	err error) {
	var conn clickhouse.Conn
	conn, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", url.Host, url.Port)},
		Auth: clickhouse.Auth{
			Database: url.Database,
			Username: url.User,
			Password: url.Password,
		},
		//Debug:           true,
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	ctx = clickhouse.Context(ctx, clickhouse.WithSettings(clickhouse.Settings{
		"max_block_size": 10,
	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
		// Some process logic.
	}))

	var result []struct {
		Item string
	}
	rawSQL := "select name as Item from system.tables where database='" + url.Database + "'"
	if err = conn.Select(ctx, &result, rawSQL); err != nil {
		return
	}

	for _, v := range result {
		items = append(items, v.Item)
	}
	return
}

//// DescribeDatasourceTablesClickHouse get a table list of type ClickHouse.
//func DescribeDatasourceTablesClickHouse(ctx context.Context, url *pbdatasource.ClickHouseURL) (items []string,
//	err error) {
//	var (
//		httpRequest  *http.Request
//		httpResponse *http.Response
//	)
//
//	client := &http.Client{Timeout: time.Second * 10}
//
//	dsn := fmt.Sprintf(
//		"http://%s:%d/?user=%s&password=%s&database=%s",
//		url.Host, url.Port, url.User, url.Password, url.Database,
//	)
//
//	rawSQL := "select name as item from system.tables where database='" + url.Database + "'"
//
//	httpRequest, err = http.NewRequest(http.MethodGet, dsn, strings.NewReader(rawSQL))
//	if err != nil {
//		return
//	}
//	httpResponse, err = client.Do(httpRequest)
//	if err != nil {
//		return
//	}
//
//	defer func() {
//		if httpResponse.Body != nil {
//			_ = httpResponse.Body.Close()
//		}
//	}()
//
//	var b []byte
//	b, err = ioutil.ReadAll(httpResponse.Body)
//	if err != nil {
//		return
//	}
//	respBody := string(b)
//
//	// FIXME: define a error code.
//	if httpResponse.StatusCode != http.StatusOK {
//		err = fmt.Errorf("%s request failed, http status code %d, message %s", dsn, httpResponse.StatusCode, respBody)
//		return
//	}
//
//	items = strings.Split(respBody, "\n")
//	return
//}

// DescribeDatasourceTablesMySQL get a table list of type MySQL.
/*func DescribeDatasourceTablesOracle(ctx context.Context, url *pbdatasource.OracleURL) (items []string, err error) {
	dsn := fmt.Sprintf("jdbc:oracle:thin:@%s:%d:%s", url.Host, url.Port, url.Database)

	var db *gorm.DB
	db, err = gorm.Open(oracle.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	defer func() {
		// close the connections.
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}()

	err = db.Raw("select table_name as item from information_schema.tables where  table_schema = ?",
		url.Database).Scan(&items).Error
	if err != nil {
		return
	}
	return
}
*/

func DescribeDatasourceTablesHive(ctx context.Context, url *pbdatasource.HiveURL) (items []string, err error) {
	conn, err := gohive.Connect(fmt.Sprintf("%s:%d", url.Host, url.Port), gohive.DefaultOptions)
	if err != nil {
		return nil, err
	}
	_, err = conn.Exec(fmt.Sprintf("use %s", url.Database))
	if err != nil {
		return nil, err
	}
	rs, err := conn.Query("show tables")
	if err != nil {
		return nil, err
	}
	var table string
	for rs.Next() {
		err = rs.Scan(&table)
		if err != nil {
			return nil, err
		}
		items = append(items, table)
	}
	defer conn.Close()
	return items, nil
}

func DescribeDatasourceTablesDB2(ctx context.Context, url *pbdatasource.DB2URL) (items []string, err error) {
	/*dsn := fmt.Sprintf("system/oracle@%s:%d/%s", url.Host, url.Port, url.Database)

	var db *gorm.DB
	db, err = gorm.Open(oracle.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	defer func() {
		// close the connections.
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}()

	err = db.Raw("select TABLE_NAME from USER_TABLES where TABLESPACE_NAME = ?", url.Database).Scan(&items).Error
	if err != nil {
		return
	}*/
	return
}

func DescribeDatasourceTablesSqlServer(ctx context.Context, url *pbdatasource.SqlServerURL) (items []string, err error) {
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", url.Host, url.Database, url.User, url.Password, url.Port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}
	stmt, err := conn.Prepare("select name from sysobjects where xtype='u'")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	var name = ""
	for rows.Next() {
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	defer rows.Close()
	return items, nil
}

func DescribeDatasourceTablesElasticSearch(ctx context.Context, url *pbdatasource.ElasticSearchURL) (items []string, err error) {
	if url.Version[0:1] == "6" {
		client6, err := elastic6.NewClient(
			elastic6.SetSniff(false),
			elastic6.SetURL(fmt.Sprintf("http://%s:%d/", url.Host, url.Port)),
			elastic6.SetBasicAuth(url.User, url.Password),
		)

		if err != nil {
			return nil, err
		}
		var indices, _ = client6.CatIndices().Do(context.Background())
		for _, v := range indices {
			items = append(items, v.Index)
		}
		return items, nil
	}
	if url.Version[0:1] == "7" {
		client7, err := elastic7.NewClient(
			elastic7.SetSniff(false),
			elastic7.SetURL(fmt.Sprintf("http://%s:%d/", url.Host, url.Port)),
			elastic7.SetBasicAuth(url.User, url.Password),
		)
		if err != nil {
			return nil, err
		}
		var indices, _ = client7.CatIndices().Do(context.Background())
		for _, v := range indices {
			items = append(items, v.Index)
		}
		return items, nil
	}
	return nil, nil

}

func DescribeDatasourceTablesMongodb(ctx context.Context, url *pbdatasource.MongoDbURL) (items []string, err error) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%d", url.Host, url.Port))
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	db := session.DB(url.Database)
	defer session.Close()
	return db.CollectionNames()
}

func DescribeDatasourceTablesKafka(ctx context.Context, url *pbdatasource.KafkaURL) (items []string, err error) {
	var brokes []string
	for _, value := range url.KafkaBrokers {
		brokes = append(brokes, fmt.Sprintf("%s:%d", value.Host, value.Port))
	}
	config := sarama.NewConfig()
	config.Version = sarama.V0_10_1_1
	client, err := sarama.NewClient(brokes, nil)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	return client.Topics()
}

func DescribeDatasourceTablesHbase(ctx context.Context, url *pbdatasource.HBaseURL) (items []string, err error) {
	hadoopMap := make(map[string]string)
	err = json.Unmarshal([]byte(url.Config), &hadoopMap)
	if err != nil {
		return nil, err
	}

	//hbase.zookeeper.quorum
	client := gohbase.NewAdminClient(hadoopMap["hbase.zookeeper.quorum"])
	// get all tables condition by namespace
	//req, err := hrpc.NewListTableNames(context.Background(),hrpc.ListNamespace("ns1"))
	// get all tables
	req, err := hrpc.NewListTableNames(context.Background())
	if err != nil {
		return nil, err
	}
	tables, _ := client.ListTableNames(req)
	for _, val := range tables {
		items = append(items, string(val.Namespace)+":"+string(val.Qualifier))
	}
	return items, nil
}

// DescribeDataSourceTables get a table list of specified data source.
func DescribeDataSourceTables(ctx context.Context, sourceType pbmodel.DataSource_Type, sourceURL *pbmodel.DataSource_URL) (
	*pbresponse.DescribeDataSourceTables, error) {
	var items []string
	var err error
	switch sourceType {
	case pbmodel.DataSource_MySQL:
		items, err = DescribeDatasourceTablesMySQL(ctx, sourceURL.Mysql)
	case pbmodel.DataSource_PostgreSQL:
		items, err = DescribeDatasourceTablesPostgresSQL(ctx, sourceURL.Postgresql)
	case pbmodel.DataSource_ClickHouse:
		items, err = DescribeDatasourceTablesClickHouse(ctx, sourceURL.Clickhouse)
	case pbmodel.DataSource_HBase:
		//todo first
		items, err = DescribeDatasourceTablesHbase(ctx, sourceURL.Hbase)
	//case pbmodel.DataSource_Ftp:
	//	// empty no need
	//case pbmodel.DataSource_HDFS:
	//	// empty no need
	case pbmodel.DataSource_SqlServer:
		items, err = DescribeDatasourceTablesSqlServer(ctx, sourceURL.Sqlserver)
	case pbmodel.DataSource_Oracle:
		// todo
		// items, err = DescribeDatasourceTablesOracle(ctx, sourceURL.Oracle)
	case pbmodel.DataSource_DB2:
		//todo
		items, err = DescribeDatasourceTablesDB2(ctx, sourceURL.Db2)
	case pbmodel.DataSource_SapHana:
		//empty no need
	case pbmodel.DataSource_Hive:
		items, err = DescribeDatasourceTablesHive(ctx, sourceURL.Hive)
	case pbmodel.DataSource_ElasticSearch:
		items, err = DescribeDatasourceTablesElasticSearch(ctx, sourceURL.ElasticSearch)
	case pbmodel.DataSource_MongoDb:
		items, err = DescribeDatasourceTablesMongodb(ctx, sourceURL.MongoDb)
	case pbmodel.DataSource_Redis:
		//empty no need
	case pbmodel.DataSource_Kafka:
		items, err = DescribeDatasourceTablesKafka(ctx, sourceURL.Kafka)
	default:
		return nil, qerror.NotSupportSourceType.Format(sourceType)
	}
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeDataSourceTables{Items: items}
	return reply, nil
}
