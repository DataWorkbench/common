package iaas_test

import (
	"context"
	"fmt"
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
	vxnet, err := iaasClient.DescribeVxnetById(ctx, "vxnet-mmwa7o8")
	require.Nil(t, err)
	_ = vxnet
	//ipNetwork := strings.Join(strings.Split(vxnet.Router.ManagerIp, ".")[:3], ".")
	//fmt.Println(ipNetwork)
}

func TestGetBalance(t *testing.T) {
	iaasClient := getIaasClient()
	balance, err := iaasClient.GetBalance(ctx, "usr-fkc8HdBQ")
	require.Nil(t, err)
	_ = balance
}

func TestAllocateVips(t *testing.T) {
	iaasClient := getIaasClient()

	vxnetId := "vxnet-mmwa7o8"
	//owner := "usr-saZdUr2m"
	//vipName := "dataomnis-flink-cfi-xxxxxxxxxxxxxxxx"
	//
	//jobId, vips, err := iaasClient.AllocateVips(ctx, &iaas.AllocateVipsInput{
	//	VxnetId:    vxnetId,
	//	VipName:    vipName,
	//	TargetUser: owner,
	//	VipAddrs:   nil,
	//	VipRange:   "172.20.0.105-172.20.0.106",
	//})
	//require.Nil(t, err)
	//fmt.Println(vips)
	//
	//// check and wait the job success.
	//for {
	//	jobSet, err := iaasClient.DescribeJobById(ctx, jobId)
	//	require.Nil(t, err)
	//	if jobSet.Status == iaas.JobSetStatusSuccessful {
	//		fmt.Println("vip create successful")
	//		break
	//	}
	//	time.Sleep(time.Second * 2)
	//}

	output, err := iaasClient.DescribeVips(ctx, &iaas.DescribeVipsInput{
		Limit:    20,
		Offset:   0,
		Vxnets:   []string{vxnetId},
		Vips:     nil,
		VipAddrs: nil,
		//VipName:  vipName,
		Owner: "",
	})
	require.Nil(t, err)

	for _, vip := range output.VipSet {
		fmt.Println(vip.VipAddr)
	}

	//vips := []string{"vip-5kn6n07j", "vip-d7bcte3p", "vip-sqco6hfb", "vip-1q3jze3a"}
	//err = iaasClient.ReleaseVips(ctx, vips)
	//require.Nil(t, err)
}

//func TestDescribeAllVxnetResources(t *testing.T)  {
//	iaasClient := getIaasClient()
//	vxnetId := "vxnet-mmwa7o8"
//
//	vxnetResourceSet, err := iaasClient.DescribeAllVxnetResources(ctx, vxnetId)
//	require.Nil(t, err)
//	fmt.Println(len(vxnetResourceSet))
//}
//
