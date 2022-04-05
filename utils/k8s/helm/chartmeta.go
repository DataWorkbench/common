package helm

import (
	"encoding/json"
	"time"
)

const DefaultTimeoutSecond = 600
const InstanceLabelKey = "app.kubernetes.io/instance"

// firstAtAll, create docker registry secret by kubectl:
// ***************************************************************
// ImageConfig
// ***************************************************************
// kubectl create secret docker-registry my-docker-registry-secret
//                                       --docker-server=<your-registry-server>
//                                       --docker-username=<your-name>
//                                       --docker-password=<your-pword>
//                                       --docker-email=<your-email>
type Image struct {
	// TODO: check pull secrets
	Registry    string   `json:"registry,omitempty" yaml:"registry,omitempty"`
	PullSecrets []string `json:"pullSecrets,omitempty" yaml:"pullSecrets,omitempty"`
	PullPolicy  string   `json:"pullPolicy,omitempty" yaml:"pullPolicy,omitempty"`

	Tag string `json:"tag,omitempty" yaml:"-"`
}

// update from other if the field is ""
func (i *Image) Update(source *Image) {
	if source == nil {
		return
	}
	if i.Registry == "" && source.Registry != "" {
		i.Registry = source.Registry
	}
	if len(i.PullSecrets) == 0 && len(source.PullSecrets) > 0 {
		i.PullSecrets = source.PullSecrets
	}
	if i.PullPolicy == "" && source.PullPolicy != "" {
		i.PullPolicy = source.PullPolicy
	}
}

type Resource struct {
	Cpu    float32 `json:"cpu,omitempty"    yaml:"cpu,omitempty"`
	Memory string `json:"memory,omitempty" yaml:"memory,omitempty"`
}

type Resources struct {
	Limits   Resource `json:"limits,omitempty"   yaml:"limits,omitempty"`
	Requests Resource `json:"requests,omitempty" yaml:"requests,omitempty"`
}

// ***************************************************************
// ValuesMeta
// ***************************************************************
type ValuesMeta struct {
	Image     *Image     `json:"image,omitempty"`
	Resources *Resources `json:"resources,omitempty"`
}

// ***************************************************************
// implement Chart interface
// ***************************************************************
type ChartMeta struct {
	// for pod label
	ChartName   string
	ReleaseName string

	Waiting       bool
	TimeOutSecond int

	DryRun bool

	Conf interface{} `json:"values,omitempty"`
}

func (m ChartMeta) ParseValues() (string, error) {
	if m.Conf != nil {
		bytes, err := json.Marshal(m.Conf)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
	return "", nil

}

func (m ChartMeta) GetChartName() string {
	return m.ChartName
}

func (m ChartMeta) GetReleaseName() string {
	return m.ReleaseName
}

func (m ChartMeta) GetLabels() map[string]string {
	return map[string]string{
		InstanceLabelKey: m.ReleaseName,
	}
}

func (m ChartMeta) WaitingReady() bool {
	return m.Waiting
}

func (m ChartMeta) GetTimeoutSecond() time.Duration {
	if m.TimeOutSecond < 1 {
		return DefaultTimeoutSecond * time.Second
	}
	return time.Duration(m.TimeOutSecond) * time.Second
}

func (m ChartMeta) IsDryRun() bool {
	return m.DryRun
}
