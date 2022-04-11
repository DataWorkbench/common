package k8s

import (
	"context"
	"github.com/DataWorkbench/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const DefaultKubeConf = "/root/.kube/config"

// **************************************************************
// the Proxy of kube client to access k8s resource
// **************************************************************
type Proxy struct {
	Client *kubernetes.Clientset
	Logger *glog.Logger
}

func (p *Proxy) GetKubeNodes(ctx context.Context) ([]string, error) {
	nodeList, err := p.Client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var nodeSlice []string
	for _, node := range nodeList.Items {
		nodeSlice = append(nodeSlice, node.Name)
	}
	return nodeSlice, nil
}

func (p *Proxy) CreateNamespace(ctx context.Context, namespace string) error {
	ns := &corev1.Namespace{}
	ns.Name = namespace
	_, err := p.Client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err = p.Client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
		}
	}
	return err
}

// Note: need to init p.kubeClient before
func (p Proxy) CheckPodsReady(ctx context.Context, namespace string, ops metav1.ListOptions) (bool, error) {
	// get PodLists
	pods, err := p.Client.CoreV1().Pods(namespace).List(ctx, ops)
	if err != nil {
		return false, err
	}
	for _, pod := range pods.Items {
		if pod.Status.Phase == corev1.PodSucceeded {
			continue
		}

		for _, condition := range pod.Status.Conditions {
			if condition.Status != corev1.ConditionTrue {
				p.Logger.Info().String("pod", pod.GetName()).
					String("is not ready, status of conditionType", string(condition.Type)).
					String("is not true, reason", condition.Reason).
					String("message", condition.Message).Fire()
				return false, nil
			}
		}
	}
	return true, nil
}

func (p Proxy) CopyConfigmap(ctx context.Context, oriNamespace, namespace, name string) error {
	_, err := p.Client.CoreV1().ConfigMaps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
		} else {
			return err
		}

		// not exist: get configmap from oriNamespace and create at namespace
		var cm *corev1.ConfigMap
		if cm, err = p.Client.CoreV1().ConfigMaps(oriNamespace).Get(ctx, name, metav1.GetOptions{}); err != nil {
			return err
		}

		// new configmap
		newCm := &corev1.ConfigMap{}
		newCm.Namespace = namespace
		newCm.Name = name
		newCm.Data = cm.Data
		newCm.BinaryData = cm.BinaryData
		_, err = p.Client.CoreV1().ConfigMaps(namespace).Create(ctx, newCm, metav1.CreateOptions{})
	}

	// exist, return
	return err
}

// if kubeConfPath == "", create k8s client auth by ServiceAccount in RBAC (/var/run/secrets/kubernetes.io/serviceaccount)
// otherwise kube client auth by kubeConfig in kubeConfPaths
func NewProxy(kubeConfPath string, logger *glog.Logger) (*Proxy, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return nil, err
	}

	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		Client: kc,
		Logger: logger,
	}, nil
}
