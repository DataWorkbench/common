package kafka

import (
	"time"

	prometheusmetrics "github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
)

// Use the same metric registry in a service.
var metricRegistry = metrics.NewRegistry()

// CollectPrometheusMetric collects sarama's metric to prometheus.
// Use `go kafka.CollectPrometheusMetric()` to make run at backend.
func CollectPrometheusMetric() {
	p := prometheusmetrics.NewPrometheusProvider(
		metricRegistry, "samara", "broker", prometheus.DefaultRegisterer, time.Second*5)
	p.UpdatePrometheusMetrics()
}
