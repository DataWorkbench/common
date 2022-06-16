package helm

import (
	"context"
	"encoding/json"
	helm "github.com/mittwald/go-helm-client"
	"time"
)

const DefaultTimeoutSecond = 30 * 60 * time.Second
const InstanceLabelKey = "app.kubernetes.io/instance"

const (
	DebugKey  = "debug"  // if enabled debug
	WaitKey   = "wait"   // if enabled wait
	DryRunKey = "dryRun" // if enabled dry-run
)

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

//func (c *Chart) GetLabels() map[string]string {
//	return map[string]string{
//		InstanceLabelKey: c.ReleaseName,
//	}
//}

// chartName: full path of HelmChart
// conf: configuration of chart, eg: from file values.yaml
func NewChartSpec(ctx context.Context, namespace, releaseName, chartName string, conf map[string]interface{}) (*helm.ChartSpec, error) {
	values, err := parseValues(conf)
	if err != nil {
		return nil, err
	}

	dryRun, ok := ctx.Value(DryRunKey).(bool)
	if !ok {
		dryRun = false
	}

	wait, ok := ctx.Value(WaitKey).(bool)
	if !ok {
		wait = false
	}

	return &helm.ChartSpec{
		Namespace:       namespace,
		CreateNamespace: true,
		ReleaseName:     releaseName,
		ChartName:       chartName,
		ValuesYaml:      values,
		Wait:            wait,
		DryRun:          dryRun,
		Timeout:         DefaultTimeoutSecond,
	}, err
}
