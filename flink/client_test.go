package flink

import (
	"context"
	"fmt"
	"github.com/DataWorkbench/glog"
	"testing"
)

var client *Client
var flinkUrl = "localhost:8081"
var ctx context.Context
var lp *glog.Logger

func init() {
	lp = glog.NewDefault().WithLevel(glog.Level(1))
	ctx = glog.WithContext(context.Background(), lp)
	client = NewClient(ctx, nil)
}

func Test_ListJobs(t *testing.T) {
	jobs, err := client.ListJobs(ctx, flinkUrl)
	if err != nil {
		t.Error(err)
	}
	if len(jobs) > 0 {
		for _, j := range jobs {
			fmt.Println(j)
		}
	}
}

func Test_GetInfo(t *testing.T) {
	info, err := client.GetInfo(ctx, flinkUrl, "7b65495367995459333e6948a37eecd4")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(info)
	}
}

func Test_GetExceptions(t *testing.T) {
	exceptions, err := client.GetExceptions(ctx, flinkUrl, "7b65495367995459333e6948a37eecd4")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(exceptions.RootException)
		fmt.Println("=================================================")
		for i, exception := range exceptions.AllExceptions {
			fmt.Println(i, exception)
			fmt.Println("=================================================")
		}
	}
}

func Test_Cancel(t *testing.T) {
	err := client.CancelJob(ctx, flinkUrl, "79216c4ca1d4a151e39f84b178e81770")
	if err != nil {
		t.Error(err)
	}
}
