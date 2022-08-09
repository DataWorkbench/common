package helm

import (
	"encoding/json"
	"github.com/DataWorkbench/common/utils/k8s"
	helm "github.com/mittwald/go-helm-client"
	"path"
	"time"
)

const DefaultTimeoutSecond = 60 * 60 * time.Second

type Config struct {
	KubeConfPath string
	HelmRepoPath string

	Debug  bool
	DryRun bool
	// if wait release ready
	WaitReady bool

	// timeout(second) of waiting release ready
	Timeout uint
}

func NewConfig(kubeConf, helmRepo string, debug, dryRun, waitReady bool, timeout uint) *Config {
	if kubeConf == "" {
		kubeConf = k8s.DefaultKubeConf
	}
	if helmRepo == "" {
		helmRepo = DefaultHelmRepoCache
	}
	return &Config{
		KubeConfPath: kubeConf,
		HelmRepoPath: helmRepo,
		Debug:        debug,
		DryRun:       dryRun,
		WaitReady:    waitReady,
		Timeout:      timeout,
	}
}

func parseValues(conf map[string]interface{}) (string, error) {
	if conf != nil {
		bytes, err := json.Marshal(conf)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return "", nil
}

// NewChartSpec
// chartName: the HelmChart filename
// valueConf: configuration of chart, eg: from file values.yaml
// conf: the optional configuration of ChartSpec, dryRun / wait / timeout(second)
func NewChartSpec(namespace, releaseName, chartName string, valueConf map[string]interface{}, conf Config) (*helm.ChartSpec, error) {
	values, err := parseValues(valueConf)
	if err != nil {
		return nil, err
	}

	var timeoutSecond = DefaultTimeoutSecond
	if conf.Timeout > 0 {
		timeoutSecond = time.Duration(conf.Timeout) * time.Second
	}

	repo := conf.HelmRepoPath
	if repo == "" {
		repo = DefaultHelmRepoCache
	}
	chart := path.Join(repo, chartName)

	return &helm.ChartSpec{
		Namespace:       namespace,
		CreateNamespace: true,
		ReleaseName:     releaseName,
		ChartName:       chart,
		ValuesYaml:      values,
		Wait:            conf.WaitReady,
		DryRun:          conf.DryRun,
		Timeout:         timeoutSecond,
	}, err
}
