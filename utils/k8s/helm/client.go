package helm

import (
	"context"
	"fmt"
	"github.com/DataWorkbench/glog"
	ghc "github.com/mittwald/go-helm-client"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/kube"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

const (
	DefaultHelmRepoConfig = ""
	DefaultHelmRepoCache  = "/root/.cache/helm/repository"

	ReleaseNotFoundErr = "release: not found"

	AllResource = "all"

	WaitInitDuration = 5 * time.Second
)

func NewClient(ctx context.Context, namespace string, conf *Config) (ghc.Client, error) {
	logger := glog.FromContext(ctx)

	debugLog := func(format string, v ...interface{}) {}
	if conf.Debug {
		debugLog = func(format string, v ...interface{}) {
			logger.Debug().Msg(fmt.Sprintf(format, v...)).Fire()
		}
	}
	opts := &ghc.Options{
		Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
		RepositoryCache:  conf.HelmRepoPath,
		RepositoryConfig: DefaultHelmRepoConfig,
		Debug:            conf.Debug,
		Linting:          true, // Change this to false if you don't want linting.
		DebugLog:         debugLog,
	}

	var restConf *rest.Config
	var err error
	if restConf, err = clientcmd.BuildConfigFromFlags("", conf.KubeConfPath); err != nil {
		return nil, err
	}
	restConfopts := &ghc.RestConfClientOptions{
		Options:    opts,
		RestConfig: restConf,
	}
	return ghc.NewClientFromRestConf(restConfopts)
}

func Exist(client ghc.Client, releaseName string) (bool, error) {
	_, err := client.GetRelease(releaseName)
	if err != nil {
		if errors.Cause(err).Error() == ReleaseNotFoundErr { // release not found
			err = nil
		}
		return false, err
	}
	return true, err
}

// WaitingResourceReady
// Any oneof labelSelector and resourceTypeAndName must be specified
func WaitingResourceReady(namespace, labelSelector string, conf Config, logFunc func(string, ...interface{}), resourceTypeAndName ...string) error {
	if len(resourceTypeAndName) == 0 {
		resourceTypeAndName = append(resourceTypeAndName, AllResource)
	}

	var restConf *rest.Config
	var err error
	var clientgetter genericclioptions.RESTClientGetter
	var client *kube.Client

	if restConf, err = clientcmd.BuildConfigFromFlags("", conf.KubeConfPath); err != nil {
		return err
	}
	clientgetter = ghc.NewRESTClientGetter(namespace, nil, restConf)
	client = kube.New(clientgetter)
	client.Log = logFunc
	builder := client.Factory.NewBuilder()

	time.Sleep(WaitInitDuration)

	result := builder.Unstructured().
		NamespaceParam(namespace).DefaultNamespace().
		LabelSelectorParam(labelSelector).
		ResourceTypeOrNameArgs(true, resourceTypeAndName...).
		ContinueOnError().
		Latest().
		Flatten().
		Do()
	infos, err := result.Infos()
	if err != nil {
		return err
	}
	timeout := DefaultTimeoutSecond
	if conf.Timeout > 0 {
		timeout = time.Duration(conf.Timeout) * time.Second
	}
	return client.Wait(infos, timeout)
}
