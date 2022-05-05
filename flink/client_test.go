package flink

import (
	"context"
	"fmt"
	"testing"

	"github.com/DataWorkbench/glog"
	"github.com/stretchr/testify/require"
)

var client *Client
var flinkUrl = "127.0.0.1:8081"
var ctx context.Context
var lp *glog.Logger

func init() {
	lp = glog.NewDefault().WithLevel(glog.Level(1))
	ctx = glog.WithContext(context.Background(), lp)
	client = New(ctx, nil)
}

func Test_Overview(t *testing.T) {
	info, err := client.Overview(ctx, flinkUrl)
	require.Nil(t, err)
	fmt.Printf("%+v\n", info)
}

func Test_ListJobs(t *testing.T) {
	jobs, err := client.ListJobs(ctx, flinkUrl)
	require.Nil(t, err)

	if len(jobs.Jobs) > 0 {
		for _, j := range jobs.Jobs {
			fmt.Printf("%+v\n", j)
		}
	}
}

func Test_DescribeJob(t *testing.T) {
	info, err := client.DescribeJob(ctx, flinkUrl, "74aa1c9bd04c33c486c0e127afd7bffa")
	require.Nil(t, err)
	fmt.Printf("%+v\n", info)
}

func Test_ListTaskManagers(t *testing.T) {
	jobs, err := client.ListTaskManagers(ctx, flinkUrl)
	require.Nil(t, err)

	if len(jobs.TaskManagers) > 0 {
		for _, j := range jobs.TaskManagers {
			fmt.Printf("%+v\n", j)
		}
	}
}

func Test_GetExceptions(t *testing.T) {
	exceptions, err := client.DescribeJobExceptions(ctx, flinkUrl, "7b65495367995459333e6948a37eecd4")
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
