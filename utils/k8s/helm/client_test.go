package helm

import (
	"fmt"
	"github.com/DataWorkbench/common/utils/k8s"
	"testing"
)

const Namespace = "dataomnis-system"
const ResourceType = "statefulsets.apps"
const ResourceName = "drc-redis-cluster-0"

const MysqlLabelSelector = "app.kubernetes.io/instance=mysql-cluster-pxc-db"
const HdfsLabelSelector = "dataomnis.io/cluster-name=hdfs-cluster"

func logFunc() func(string, ...interface{}) {
	return func(s string, i ...interface{}) {
		println(fmt.Sprintf(s, i...))
	}
}

func TestWaitingResource(t *testing.T) {
	t.Logf("waitting with resourceType and resourceName ..")
	err := WaitingResourceReady(Namespace, k8s.DefaultKubeConf, "", DefaultTimeoutSecond, logFunc(), ResourceType, ResourceName)
	if err != nil {
		t.Fatalf("waiting resource(%s/%s) ready error: %v", ResourceType, ResourceName, err)
	}

	t.Logf("waitting with labelSelector ..")
	err = WaitingResourceReady(Namespace, k8s.DefaultKubeConf, HdfsLabelSelector, DefaultTimeoutSecond, logFunc())
	if err != nil {
		t.Fatalf("waiting resource(%s/%s) ready error: %v", ResourceType, ResourceName, err)
	}
}
