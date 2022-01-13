package zeppelin

//
//import (
//	"fmt"
//	"testing"
//	"time"
//)
//
//var config = ClientConfig{
//	ZeppelinRestUrl: "http://127.0.0.1:8080",
//	Timeout:         time.Millisecond * 2000,
//	RetryCount:      2,
//	QueryInterval:   2000,
//}
//
//var zSession *ZSession
//
//func init() {
//	var properties = map[string]string{}
//	properties["FLINK_HOME"] = "/Users/apple/develop/bigdata/flink-1.12.5"
//	properties["flink.execution.mode"] = "remote"
//	properties["flink.execution.remote.host"] = "localhost"
//	properties["flink.execution.remote.port"] = "8081"
//	zSession = NewZSessionWithProperties(config, "flink", properties)
//}
//
//func Test_Start(t *testing.T) {
//	err := zSession.Start()
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func Test_RunSql(t *testing.T) {
//	err := zSession.Start()
//	if err != nil {
//		t.Error(err)
//	}
//	result, err := zSession.Execute("ssql", "drop table if exists datagen;")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(result)
//	result, err = zSession.Execute("ssql", "create table datagen(id int,name string) with ('connector' = 'datagen',"+
//		"'rows-per-second' = '2');")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(result)
//	result, err = zSession.Execute("ssql", "drop table if exists print;")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(result)
//	result, err = zSession.Execute("ssql", "create table print(id int,name string) with ('connector'='print');")
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(result)
//	var properties = map[string]string{}
//	properties["parallelism"] = "1"
//	properties["jobName"] = "demo01"
//	result, err = zSession.SubmitWithProperties("ssql", properties, "insert into print select * from datagen;")
//	if err != nil {
//		t.Error(err)
//	}
//	result, err = zSession.WaitUntilRunning(result.StatementId)
//	if err != nil {
//		t.Error(err)
//	}
//	fmt.Println(result)
//}
