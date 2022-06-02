package helm

import (
	"context"
	"fmt"
	"github.com/DataWorkbench/glog"
	helm "github.com/mittwald/go-helm-client"
	"github.com/pkg/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	DefaultDurationSec = 20

	DefaultHelmRepoConfig = ""
	DefaultHelmRepoCache  = "/root/.cache/helm/repository"

	ReleaseNotFoundErr = "release: not found"
)

// ******************************************************************
// helm client Proxy to handle helm release, implement Helm interface
// ******************************************************************
type Proxy struct {
	client          helm.Client // helm client
	kubeConfPath    string
	namespace       string
	repositoryCache string
	logger          *glog.Logger
}

func NewProxy(ctx context.Context, namespace, kubeConfPath string) (*Proxy, error) {
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
	opts := &helm.Options{
		Namespace:        namespace, // Change this to the namespace you wish to install the chart in.
		RepositoryCache:  DefaultHelmRepoCache,
		RepositoryConfig: DefaultHelmRepoConfig,
		Debug:            debug,
		Linting:          true, // Change this to false if you don't want linting.
		DebugLog:         debugLog,
	}

	var kubeConf *rest.Config
	var client helm.Client
	var err error
	if kubeConf, err = clientcmd.BuildConfigFromFlags("", kubeConfPath); err != nil {
		return nil, err
	}
	restConfopts := &helm.RestConfClientOptions{
		Options:    opts,
		RestConfig: kubeConf,
	}
	if client, err = helm.NewClientFromRestConf(restConfopts); err != nil {
		return nil, err
	}

	return &Proxy{
		kubeConfPath:    kubeConfPath,
		namespace:       namespace,
		repositoryCache: DefaultHelmRepoCache,
		logger:          logger,
		client:          client,
	}, err
}


func (p Proxy) InstallOrUpgrade(ctx context.Context, chart *helm.ChartSpec) error {
	p.logger.Info().String("helm install release", chart.ReleaseName).String("with chart", chart.ChartName).Fire()
	_, err := p.client.InstallOrUpgradeChart(ctx, chart)
	if err != nil {
		p.logger.Error().Error("helm install error", err).Fire()
		return err
	}
	return err
}

//func (p Proxy) WaitingReady(ctx context.Context, chart Chart) error {
//	name := chart.GetReleaseName()
//	p.logger.Info().String("waiting release", name).Msg("ready..").Fire()
//
//	labelMap := chart.GetLabels()
//	ops := v1.ListOptions{
//		LabelSelector: labels.SelectorFromSet(labelMap).String(),
//	}
//
//	ready := false
//	kProxy, err := k8s.NewProxy(p.kubeConfPath, p.logger)
//	if err != nil {
//		p.logger.Error().Error("new kube-client proxy error", err).Fire()
//		return err
//	}
//
//	duration := time.Duration(DefaultDurationSec) * time.Second
//	for {
//		select {
//		case <-time.After(duration):
//			ready, err = kProxy.CheckPodsReady(ctx, p.namespace, ops)
//			if err != nil {
//				p.logger.Error().Error("check status ready error", err).Fire()
//				return err
//			}
//			if ready {
//				p.logger.Info().String("all pods ready of release", name).
//					String("in namespace", p.namespace).Fire()
//				return nil
//			}
//		case <-ctx.Done():
//			p.logger.Warn().Error("waiting-action been canceled, error", ctx.Err()).Fire()
//			return errors.Errorf("install release=%s timeout", chart.GetReleaseName())
//		}
//	}
//}

func (p Proxy) Exist(releaseName string) (bool, error) {
	_, err := p.client.GetRelease(releaseName)
	if err != nil {
		if errors.Cause(err).Error() == ReleaseNotFoundErr { // release not found
			err = nil
		}
		return false, err
	}
	return true, err
}

func (p Proxy) Delete(releaseName string) error {
	spec := &helm.ChartSpec{
		ReleaseName:  releaseName,
		MaxHistory:   0,
		DisableHooks: true,
	}
	return p.client.UninstallRelease(spec)
}
