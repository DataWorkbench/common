package iaas_test

import (
	"context"
	"testing"

	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/DataWorkbench/glog"
	"github.com/stretchr/testify/require"
)

func loadConfig() *iaas.Config {
	cfg := &iaas.Config{
		Zone:            "testing",
		Host:            "api.testing.com",
		Port:            7777,
		Protocol:        "http",
		Timeout:         30,
		AccessKeyId:     "LTMJGBXPHSEZRNVKKPHU",
		SecretAccessKey: "7GvVuGAx2iB8NA9n8NtczH8BJnTkDGwGm9N6DYBo",
		Uri:             "/iaas/",
	}
	return cfg
}

var ctx = glog.WithContext(context.Background(), glog.NewDefault().WithLevel(glog.DebugLevel))

func getIaasClient() *iaas.Client {
	cfg := loadConfig()
	iaasClient := iaas.New(ctx, cfg)
	return iaasClient
}

func TestDescribeIaasVxnet(t *testing.T) {
	iaasClient := getIaasClient()
	vxnet, err := iaasClient.DescribeVxnetById(ctx, "vxnet-xxxxx")
	require.NotNil(t, err)
	_ = vxnet
}

func TestGetBalance(t *testing.T) {
	iaasClient := getIaasClient()
	balance, err := iaasClient.GetBalance(ctx, "usr-fkc8HdBQ")
	require.Nil(t, err)
	_ = balance
}
