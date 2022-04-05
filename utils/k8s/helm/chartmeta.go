package helm

import (
	"encoding/json"
	"time"
)

const DefaultTimeoutSecond = 600
const InstanceLabelKey = "app.kubernetes.io/instance"


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
