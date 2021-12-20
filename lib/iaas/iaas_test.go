package iaas_test

import (
	"context"
	"testing"

	"github.com/DataWorkbench/common/lib/iaas"
	"github.com/DataWorkbench/glog"
)

func TestDescribeIaasVxnet(t *testing.T) {
	cfg := &iaas.Config{
		Zone:            "testing1a",
		Host:            "api.testing.com",
		Port:            7777,
		Protocol:        "http",
		Timeout:         30,
		AccessKeyId:     "VEEZBXXQRPJMGDWAGOLY",
		SecretAccessKey: "xNwZ2ioxPvx6efwQnowtnnMwhZ3kiFzHCnpkdmKW",
	}

	ctx := glog.WithContext(context.Background(), glog.NewDefault().WithLevel(glog.DebugLevel))
	iaasClient := iaas.New(ctx, cfg)
	vxnet, err := iaasClient.DescribeVxnetById(ctx, "vxnet-kwgjf73z")
	if err != nil {
		t.Fatalf("describe vxnet error: %+v", err)
	} else {
		t.Logf("vxnet: %+v", vxnet)
	}
}
