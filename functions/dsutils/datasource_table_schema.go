package dsutils

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DataWorkbench/common/qerror"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel/pbdatasource"
	"github.com/DataWorkbench/gproto/xgo/types/pbresponse"
	"github.com/dazheng/gohive"
	_ "github.com/denisenkom/go-mssqldb"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"time"
)

func escapeColumnType(columnType string) string {
	if columnType == "" {
		return columnType
	}
	columnType = strings.Split(columnType, " ")[0]
	index := strings.Index(columnType, "(")
	if index != -1 {
		columnType = columnType[:index]
	}
	return strings.ToUpper(columnType)
}

func DescribeDatasourceTableSchemaMySQL(ctx context.Context, url *pbdatasource.MySQLURL,
	tableName string) (columns []*pbdatasource.TableColumn, err error) {
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

	var rawSQL strings.Builder
	rawSQL.Grow(512)

	rawSQL.WriteString("SELECT ")
	rawSQL.WriteString("COLUMN_NAME AS name,")
	rawSQL.WriteString("COLUMN_TYPE AS type,")
	//rawSQL.WriteString("'' as length,")
	rawSQL.WriteString("CASE COLUMN_KEY when 'PRI' then 'true' else 'false' end AS is_primary_key")
	rawSQL.WriteString(" FROM information_schema.columns")
	rawSQL.WriteString(" WHERE ")
	rawSQL.WriteString("table_schema = '" + url.Database + "'")
	rawSQL.WriteString(" AND ")
	rawSQL.WriteString("table_name = '" + tableName + "'")
	rawSQL.WriteString(";")

	columns = make([]*pbdatasource.TableColumn, 0)
	err = db.Raw(rawSQL.String()).Scan(&columns).Error
	if err != nil {
		return
	}
	for _, column := range columns {
		column.Type = escapeColumnType(column.Type)
	}

	return
}

func DescribeDatasourceTableSchemaPostgreSQL(ctx context.Context, url *pbdatasource.PostgreSQLURL,
	tableName string) (columns []*pbdatasource.TableColumn, err error) {
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

	//err = db.Raw(
	//	"SELECT a.attname as name, pg_catalog.format_type(a.atttypid, a.atttypmod) as type, '' as length, " +
	//		"case attnotnull when 't' then 'true' else 'false' end  as is_primary_key FROM pg_catalog." +
	//		"pg_attribute a WHERE a.attrelid = '" + tableName + "'::regclass::oid AND a.attnum > 0 AND NOT a.attisdropped").
	//	Scan(&resp.Columns).Error

	var rawSQL strings.Builder
	rawSQL.Grow(512)
	rawSQL.WriteString("SELECT ")
	rawSQL.WriteString("a.attname as name,")
	rawSQL.WriteString("pg_catalog.format_type(a.atttypid, a.atttypmod) as type,")
	//rawSQL.WriteString("'' as length,")
	rawSQL.WriteString("case attnotnull when 't' then 'true' else 'false' end  as is_primary_key")
	rawSQL.WriteString(" FROM pg_catalog.pg_attribute a")
	rawSQL.WriteString(" WHERE ")
	rawSQL.WriteString("a.attrelid = '" + tableName + "'::regclass::oid")
	rawSQL.WriteString(" AND ")
	rawSQL.WriteString("a.attnum > 0 ")
	rawSQL.WriteString(" AND ")
	rawSQL.WriteString("NOT a.attisdropped")
	rawSQL.WriteString(";")

	columns = make([]*pbdatasource.TableColumn, 0)
	err = db.Raw(rawSQL.String()).Scan(&columns).Error
	if err != nil {
		return
	}
	for _, column := range columns {
		column.Type = escapeColumnType(column.Type)
	}
	return
}

func DescribeDatasourceTableSchemaClickHouse(ctx context.Context, url *pbdatasource.ClickHouseURL,
	tableName string) (columns []*pbdatasource.TableColumn, err error) {
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

	var rawSQL strings.Builder
	rawSQL.Grow(512)
	rawSQL.WriteString("SELECT ")
	rawSQL.WriteString("name as Name, type as Type, is_in_primary_key as IsPrimaryKey")
	//rawSQL.WriteString("name, type, '' length, ")
	//rawSQL.WriteString("case is_in_primary_key when 1 then 'true' else 'false' end as IsPrimaryKey")
	rawSQL.WriteString(" FROM system.columns")
	rawSQL.WriteString(" WHERE ")
	rawSQL.WriteString(" table = '" + tableName + "'")
	rawSQL.WriteString(" AND ")
	rawSQL.WriteString(" database = '" + url.Database + "'")
	rawSQL.WriteString(";")

	//reqBody := strings.NewReader("select name, type, '' length, " +
	//	"case is_in_primary_key when 1 then 'true' else 'false' end as is_primary_key from system." +
	//	"columns where table= '" + tableName + "' and database='" + url.GetDatabase() + "'")

	var result []struct {
		Name         string
		Type         string
		IsPrimaryKey uint8
	}

	if err = conn.Select(ctx, &result, rawSQL.String()); err != nil {
		return
	}
	for i := 0; i < len(result); i++ {
		v := result[i]
		var isPrimaryKey bool
		if v.IsPrimaryKey == 1 {
			isPrimaryKey = true
		}
		columns = append(columns, &pbdatasource.TableColumn{
			Name:         v.Name,
			Type:         escapeColumnType(v.Type),
			IsPrimaryKey: isPrimaryKey,
		})
	}
	return
}

func DescribeDatasourceTableSchemaHbase(ctx context.Context, url *pbdatasource.HBaseURL,
	tableName string) (columns []*pbdatasource.TableColumn, err error) {

	return
}

//func DescribeDatasourceTableSchemaClickHouse(ctx context.Context, url *pbdatasource.ClickHouseURL,
//	tableName string) (columns []*pbdatasource.TableColumn, err error) {
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
//	var rawSQL strings.Builder
//	rawSQL.Grow(512)
//	rawSQL.WriteString("SELECT ")
//	rawSQL.WriteString("name, type, ")
//	//rawSQL.WriteString("name, type, '' length, ")
//	rawSQL.WriteString("case is_in_primary_key when 1 then 'true' else 'false' end as is_primary_key")
//	rawSQL.WriteString(" FROM system.columns")
//	rawSQL.WriteString(" WHERE ")
//	rawSQL.WriteString(" table = '" + tableName + "'")
//	rawSQL.WriteString(" AND ")
//	rawSQL.WriteString(" database = '" + url.Database + "'")
//	rawSQL.WriteString(";")
//
//	//reqBody := strings.NewReader("select name, type, '' length, " +
//	//	"case is_in_primary_key when 1 then 'true' else 'false' end as is_primary_key from system." +
//	//	"columns where table= '" + tableName + "' and database='" + url.GetDatabase() + "'")
//
//	httpRequest, err = http.NewRequest(http.MethodGet, dsn, strings.NewReader(rawSQL.String()))
//	if err != nil {
//		return
//	}
//
//	httpResponse, err = client.Do(httpRequest)
//	if err != nil {
//		return
//	}
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
//
//	respBody := string(b)
//	if httpResponse.StatusCode != http.StatusOK {
//		err = fmt.Errorf("%s request failed, http status code %d, message %s", dsn, httpResponse.StatusCode, respBody)
//		return
//	}
//
//	columns = make([]*pbdatasource.TableColumn, 0)
//
//	returnItems := strings.Split(respBody, "\n")
//	for i := 0; i < len(returnItems); i++ {
//		dbColumn := strings.Split(returnItems[i], "	")
//
//		columnName := dbColumn[0]
//		columnType := dbColumn[1]
//		isPrimaryKey, xErr := strconv.ParseBool(dbColumn[2])
//		if xErr != nil {
//			return nil, xErr
//		}
//		column := &pbdatasource.TableColumn{
//			Name:         columnName,
//			Type:         columnType,
//			IsPrimaryKey: isPrimaryKey,
//		}
//		columns = append(columns, column)
//	}
//	return
//}

func DescribeDatasourceTableSchemaSqlServer(ctx context.Context, url *pbdatasource.SqlServerURL,
	tableName string) (columns []*pbdatasource.TableColumn, err error) {
	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", url.Host, url.Database, url.User, url.Password, url.Port)
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		return nil, err
	}

	var sqlStr = `SELECT
        col.name AS name ,
        t.name AS type ,
        CASE WHEN EXISTS ( SELECT   1
                           FROM   dbo.sysindexes si
                               INNER JOIN dbo.sysindexkeys sik ON si.id = sik.id
                               AND si.indid = sik.indid
                                        INNER JOIN dbo.syscolumns sc ON sc.id = sik.id
                               AND sc.colid = sik.colid
                                        INNER JOIN dbo.sysobjects so ON so.name = si.name
                               AND so.xtype = 'PK'
                           WHERE    sc.id = col.id
                             AND sc.colid = col.colid ) THEN '1'
					 ELSE ''
					END AS is_primary_key
		FROM  dbo.syscolumns col
					LEFT  JOIN dbo.systypes t ON col.xtype = t.xusertype
					inner JOIN dbo.sysobjects obj ON col.id = obj.id
			AND obj.xtype = 'U'
			AND obj.status >= 0
					LEFT  JOIN dbo.syscomments comm ON col.cdefault = comm.id
					LEFT  JOIN sys.extended_properties ep ON col.id = ep.major_id
			AND col.colid = ep.minor_id
			AND ep.name = 'MS_Description'
					LEFT  JOIN sys.extended_properties epTwo ON obj.id = epTwo.major_id
			AND epTwo.minor_id = 0
			AND epTwo.name = 'MS_Description'
		WHERE obj.name = '%s'
		ORDER BY col.colorder`

	stmt, err := conn.Prepare(fmt.Sprintf(sqlStr, tableName))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
		var columnName string
		var columnType string
		var primaryKey uint8
		rows.Scan(&columnName, &columnType, &primaryKey)
		var isPrimaryKey bool
		if primaryKey == 1 {
			isPrimaryKey = true
		}
		columns = append(columns, &pbdatasource.TableColumn{
			Name:         columnName,
			Type:         columnType,
			IsPrimaryKey: isPrimaryKey,
		})
	}
	return
}

func DescribeDatasourceTableSchemaHive(ctx context.Context, url *pbdatasource.HiveURL,
	tableName string) (columns []*pbdatasource.TableColumn, err error) {
	conn, err := gohive.Connect(fmt.Sprintf("%s:%d", url.Host, url.Port), gohive.DefaultOptions)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(fmt.Sprintf("desc %s", tableName))
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	for rows.Next() {
		var columnName string
		var columnType string
		var comment string
		rows.Scan(&columnName, &columnType, &comment)
		columns = append(columns, &pbdatasource.TableColumn{
			Name:         columnName,
			Type:         columnType,
			IsPrimaryKey: false,
		})
	}
	return
}

//DescribeDataSourceTableSchema get the table schema of specified table in datasource.
func DescribeDataSourceTableSchema(ctx context.Context, sourceType pbmodel.DataSource_Type, sourceURL *pbmodel.DataSource_URL, tableName string) (*pbresponse.DescribeDataSourceTableSchema, error) {

	var columns []*pbdatasource.TableColumn
	var err error
	switch sourceType {
	case pbmodel.DataSource_MySQL:
		columns, err = DescribeDatasourceTableSchemaMySQL(ctx, sourceURL.Mysql, tableName)
	case pbmodel.DataSource_PostgreSQL:
		columns, err = DescribeDatasourceTableSchemaPostgreSQL(ctx, sourceURL.Postgresql, tableName)
	case pbmodel.DataSource_ClickHouse:
		columns, err = DescribeDatasourceTableSchemaClickHouse(ctx, sourceURL.Clickhouse, tableName)
	case pbmodel.DataSource_SqlServer:
		columns, err = DescribeDatasourceTableSchemaSqlServer(ctx, sourceURL.Sqlserver, tableName)
	case pbmodel.DataSource_Oracle:
		//empty
	case pbmodel.DataSource_DB2:
		//empty
	case pbmodel.DataSource_SapHana:
		//empty
	case pbmodel.DataSource_Hive:
		columns, err = DescribeDatasourceTableSchemaHive(ctx, sourceURL.Hive, tableName)
	case pbmodel.DataSource_HBase:
		columns, err = DescribeDatasourceTableSchemaHbase(ctx, sourceURL.Hbase, tableName)

	default:
		return nil, qerror.NotSupportSourceType.Format(sourceType)
	}
	if err != nil {
		return nil, err
	}
	reply := &pbresponse.DescribeDataSourceTableSchema{Schema: &pbdatasource.TableSchema{
		Columns: columns,
	}}
	return reply, nil
}
