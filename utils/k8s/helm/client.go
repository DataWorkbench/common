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
)


func NewClient(ctx context.Context, namespace, kubeConfPath string) (ghc.Client, error) {
	logger := glog.FromContext(ctx)
	debug, ok := ctx.Value(DebugKey).(bool)
	if !ok {
		debug = false
	}

	debugLog := func(format string, v ...interface{}) {}
	if debug {
		debugLog = func(format string, v ...interface{}) {
			logger.Debug().Msg(fmt.Sprintf(format, v...)).Fire()
		}
	}
	opts := &ghc.Options{
		Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
		RepositoryCache:  DefaultHelmRepoCache,
		RepositoryConfig: DefaultHelmRepoConfig,
		Debug:            debug,
		Linting:          true, // Change this to false if you don't want linting.
		DebugLog:         debugLog,
	}

	var conf *rest.Config
	var err error
	if conf, err = clientcmd.BuildConfigFromFlags("", kubeConfPath); err != nil {
		return nil, err
	}
	restConfopts := &ghc.RestConfClientOptions{
		Options:    opts,
		RestConfig: conf,
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


// Any oneof labelSelector and resourceTypeAndName must be specified
func WaitingResourceReady(namespace, kubeConfPath, labelSelector string, timeout time.Duration, logFunc func(string, ...interface{}), resourceTypeAndName ...string) error {
	if len(resourceTypeAndName) == 0 {
		resourceTypeAndName = append(resourceTypeAndName, AllResource)
	}

	var conf *rest.Config
	var err error
	var clientgetter genericclioptions.RESTClientGetter
	var client *kube.Client

	if conf, err = clientcmd.BuildConfigFromFlags("", kubeConfPath); err != nil {
		return err
	}
	clientgetter = ghc.NewRESTClientGetter(namespace, nil, conf)
	client = kube.New(clientgetter)
	client.Log = logFunc
	builder := client.Factory.NewBuilder()

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
	return client.Wait(infos, timeout)
}
