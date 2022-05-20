package dsutils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dazheng/gohive"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DataWorkbench/gproto/xgo/types/pbmodel"
	"github.com/DataWorkbench/gproto/xgo/types/pbmodel/pbdatasource"
	"github.com/Shopify/sarama"
	"github.com/dutchcoders/goftp"
	elastic6 "github.com/olivere/elastic/v6"
	elastic7 "github.com/olivere/elastic/v7"
	"github.com/samuel/go-zookeeper/zk"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func pingMysql(url *pbdatasource.MySQLURL) (err error) {
	var conn net.Conn
	ip := net.JoinHostPort(url.Host, strconv.Itoa(int(url.Port)))
	conn, err = net.DialTimeout("tcp", ip, time.Second*3)
	if err != nil {
		return
	}
	if conn != nil {
		_ = conn.Close()
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		url.User, url.Password, url.Host, url.Port, url.Database,
	)
	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	} else {
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}
	return
}

func pingPostgreSQL(url *pbdatasource.PostgreSQLURL) (err error) {
	var conn net.Conn

	ip := net.JoinHostPort(url.Host, strconv.Itoa(int(url.Port)))
	conn, err = net.DialTimeout("tcp", ip, time.Second*3)
	if err != nil {
		return
	}
	if conn != nil {
		_ = conn.Close()
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d  dbname=%s ",
		url.User, url.Password, url.Host, url.Port, url.Database,
	)
	var db *gorm.DB
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	} else {
		if sqlDB, e := db.DB(); e == nil {
			_ = sqlDB.Close()
		}
	}
	return
}

func pingSqlServer(url *pbdatasource.SqlServerURL) (err error) {
	var conn net.Conn

	ip := net.JoinHostPort(url.Host, strconv.Itoa(int(url.Port)))
	conn, err = net.DialTimeout("tcp", ip, time.Second*3)
	if err != nil {
		return
	}
	if conn != nil {
		_ = conn.Close()
	}

	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d;encrypt=disable", url.Host, url.Database, url.User, url.Password, url.Port)
	_, err = sql.Open("mssql", connString)
	if err != nil {
		return err
	}
	return
}

func pingClickHouse(url *pbdatasource.ClickHouseURL) (err error) {
	var conn net.Conn
	ip := net.JoinHostPort(url.Host, strconv.Itoa(int(url.Port)))
	conn, err = net.DialTimeout("tcp", ip, time.Second*3)
	if err != nil {
		return
	}
	if conn != nil {
		_ = conn.Close()
	}

	var (
		client  *http.Client
		req     *http.Request
		rep     *http.Response
		reqBody io.Reader
	)

	client = &http.Client{Timeout: time.Millisecond * 100}
	reqBody = strings.NewReader("SELECT 1")
	dsn := fmt.Sprintf(
		"http://%s:%d/?user=%s&password=%s&database=%s",
		url.Host, url.Port, url.User, url.Password, url.Database,
	)

	req, err = http.NewRequest(http.MethodGet, dsn, reqBody)
	if err != nil {
		return
	}

	rep, err = client.Do(req)
	if err != nil {
		return
	}

	repBody, _ := ioutil.ReadAll(rep.Body)
	_ = rep.Body.Close()

	if rep.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s request failed, http status code %d, message %s", dsn, rep.StatusCode, string(repBody))
		_ = rep.Body.Close()
		return
	}
	return
}

func pingKafka(url *pbdatasource.KafkaURL) (err error) {
	dsn := make([]string, 0, len(url.KafkaBrokers))

	for _, item := range url.KafkaBrokers {
		dsn = append(dsn, fmt.Sprintf("%s:%d", item.Host, item.Port))
	}

	consumer, terr := sarama.NewConsumer(dsn, nil)
	if terr != nil {
		err = terr
		return
	}
	_ = consumer.Close()
	return
}

func pingHBase(url *pbdatasource.HBaseURL) (err error) {
	var conn *zk.Conn

	config := make(map[string]string)
	if err = json.Unmarshal([]byte(url.Config), &config); err != nil {
		return err
	}

	zkHosts := config["hbase.zookeeper.quorum"]
	zkPort := config["hbase.zookeeper.property.clientPort"]
	if zkPort == "" {
		zkPort = "2181"
	}

	servers := make([]string, 0, len(zkHosts))
	for _, node := range strings.Split(zkHosts, ",") {
		servers = append(servers, fmt.Sprintf("%s:%s", node, zkPort))
	}

	conn, _, err = zk.Connect(servers, time.Millisecond*100)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func pingFtp(url *pbdatasource.FtpURL) (err error) {
	var (
		conn *goftp.FTP
	)
	if conn, err = goftp.Connect(fmt.Sprintf("%v:%d", url.Host, url.Port)); err != nil {
		return err
	}
	err = conn.Login(url.User, url.Password)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func pingHDFS(url *pbdatasource.HDFSURL) (err error) {
	var conn net.Conn
	// https://github.com/colinmarc/hdfs -- install the hadoop client. so don't use it.
	// https://studygolang.com/articles/766 -- use 50070 http port. but user input the IPC port.
	ip := net.JoinHostPort(url.NameNode, strconv.Itoa(int(url.Port)))
	conn, err = net.DialTimeout("tcp", ip, time.Millisecond*2000)
	if err != nil {
		return err
	}
	if conn != nil {
		_ = conn.Close()
	}
	return err
}

func pingHive(url *pbdatasource.HiveURL) (err error) {
	conn, err := gohive.Connect(fmt.Sprintf("%s:%d", url.Host, url.Port), gohive.DefaultOptions)
	if err != nil {
		return err
	}
	defer conn.Close()
	return nil
}

func pingElasticSearch(url *pbdatasource.ElasticSearchURL) (err error) {
	if url.Version[0:1] == "6" {
		_, err := elastic6.NewClient(
			elastic6.SetSniff(false),
			elastic6.SetURL(fmt.Sprintf("http://%s:%d/", url.Host, url.Port)),
			elastic6.SetBasicAuth(url.User, url.Password),
		)

		if err != nil {
			return err
		}
	}
	if url.Version[0:1] == "7" {
		_, err := elastic7.NewClient(
			elastic7.SetSniff(false),
			elastic7.SetURL(fmt.Sprintf("http://%s:%d/", url.Host, url.Port)),
			elastic7.SetBasicAuth(url.User, url.Password),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func pingMongoDb(url *pbdatasource.MongoDbURL) (err error) {
	session, err := mgo.Dial(fmt.Sprintf("%s:%d", url.Hosts[0].Host, url.Hosts[0].Port))
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	if url.User != "" && url.Password != "" {
		db := session.DB("admin")
		err = db.Login(url.User, url.Password)
		if err != nil {
			return err
		}
	}
	defer session.Close()
	return nil
}

func pingRedis(url *pbdatasource.RedisURL) (err error) {
	var redisOption = redis.Options{
		Addr: fmt.Sprintf("%s:%d", url.Hosts[0].Host, url.Hosts[0].Port),
	}
	if url.Password != "" {
		redisOption.Password = url.Password
	}

	client := redis.NewClient(&redisOption)
	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func PingDataSourceConnection(ctx context.Context, sourceType pbmodel.DataSource_Type, sourceURL *pbmodel.DataSource_URL) (connInfo *pbmodel.DataSourceConnection, err error) {
	begin := time.Now()
	message := ""
	result := pbmodel.DataSourceConnection_Success

	switch sourceType {
	case pbmodel.DataSource_MySQL:
		err = pingMysql(sourceURL.Mysql)
	case pbmodel.DataSource_PostgreSQL:
		err = pingPostgreSQL(sourceURL.Postgresql)
	case pbmodel.DataSource_SqlServer:
		err = pingSqlServer(sourceURL.Sqlserver)
	case pbmodel.DataSource_Oracle:
		//empty
	case pbmodel.DataSource_DB2:
		//empty
	case pbmodel.DataSource_SapHana:
		//empty
	case pbmodel.DataSource_Kafka:
		err = pingKafka(sourceURL.Kafka)
	case pbmodel.DataSource_S3:
		//empty
	case pbmodel.DataSource_ClickHouse:
		err = pingClickHouse(sourceURL.Clickhouse)
	case pbmodel.DataSource_HBase:
		err = pingHBase(sourceURL.Hbase)
	case pbmodel.DataSource_Ftp:
		err = pingFtp(sourceURL.Ftp)
	case pbmodel.DataSource_HDFS:
		err = pingHDFS(sourceURL.Hdfs)
	case pbmodel.DataSource_Hive:
		err = pingHive(sourceURL.Hive)
	case pbmodel.DataSource_ElasticSearch:
		err = pingElasticSearch(sourceURL.ElasticSearch)
	case pbmodel.DataSource_MongoDb:
		err = pingMongoDb(sourceURL.MongoDb)
	case pbmodel.DataSource_Redis:
		err = pingRedis(sourceURL.Redis)
	}

	if err != nil {
		result = pbmodel.DataSourceConnection_Failed
		message = err.Error()
		err = nil
	}

	connInfo = &pbmodel.DataSourceConnection{
		SpaceId:     "",
		SourceId:    "",
		NetworkId:   "",
		Status:      pbmodel.DataSourceConnection_Enabled,
		Result:      result,
		Message:     message,
		Created:     begin.Unix(),
		Elapse:      time.Since(begin).Milliseconds(),
		NetworkInfo: nil,
	}
	return
}
