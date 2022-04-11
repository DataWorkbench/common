package helm

import (
	"context"
	"time"
)

// the Values for helm release
type Values interface {
	Parse() (string, error)
}

// helm chart interface
// Chart is the proxy of helm chart that with the Values Configuration.
type Chart interface {
	// parse the field-values to Values for helm release
	ParseValues() (string, error)

	// return chart name
	GetChartName() string

	// return relase name
	GetReleaseName() string

	GetLabels() map[string]string

	// whether to wait release ready
	WaitingReady() bool
	GetTimeoutSecond() time.Duration

	IsDryRun() bool
}

// helm client interface for dataomnis-service
type Helm interface {
	// install Chart to k8s as a Release
	InstallOrUpgrade(context.Context, Chart) error

	// waiting a release ready
	WaitingReady(context.Context, Chart) error

	// check if a release by name exist
	Exist(string) (bool, error)

	// delete a release by name(string)
	Delete(string) error
}
