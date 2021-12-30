package zeppelin

import (
	"testing"
	"time"
)

func Test_Start(t *testing.T) {
	config := ClientConfig{
		ZeppelinRestUrl: "http://127.0.0.1:8080",
		Timeout:         time.Millisecond * 2000,
		RetryCount:      2,
		QueryInterval:   2000,
	}
	var properties = map[string]string{}
	properties["FLINK_HOME"] = "/Users/apple/develop/bigdata/flink-1.12.5"
	properties["flink.execution.mode"] = "remote"
	properties["flink.execution.remote.host"] = "localhost"
	properties["flink.execution.remote.port"] = "8081"
	zSession := NewZSession4(config, "flink", properties, 100)
	err := zSession.start()
	if err != nil {
		t.Error(err)
	}
}
