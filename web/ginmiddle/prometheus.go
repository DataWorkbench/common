//
// References https://github.com/zsais/go-gin-prometheus
//
package ginmiddle

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	subsystem  = "http"
	metricPath = "/metrics"
)

func Prometheus(engine *gin.Engine) gin.HandlerFunc {
	// Register metrics path to engine.
	engine.GET(metricPath, prometheusHandler())

	// Init collector
	reqCounterVec := metricReqCounterVec()
	reqDurHistVec := metricReqDurHistogramVec()
	reqSizeSummary := metricReqSizeSummary()
	respSizeSummary := metricRespSizeSummary()

	prometheus.MustRegister(reqCounterVec)
	prometheus.MustRegister(reqDurHistVec)
	prometheus.MustRegister(reqSizeSummary)
	prometheus.MustRegister(respSizeSummary)

	return func(c *gin.Context) {
		if c.Request.URL.Path == metricPath {
			c.Next()
			return
		}

		start := time.Now()
		reqSize := computeApproximateRequestSize(c.Request)

		// Call the next handler
		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		respSize := float64(c.Writer.Size())

		uri := uriLabelMapping(c)

		reqCounterVec.WithLabelValues(status, c.Request.Method, c.HandlerName(), c.Request.Host, uri).Inc()
		reqDurHistVec.WithLabelValues(status, c.Request.Method, uri).Observe(elapsed)

		reqSizeSummary.Observe(reqSize)
		respSizeSummary.Observe(respSize)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func uriLabelMapping(c *gin.Context) string {
	uri := c.Request.URL.Path
	for _, p := range c.Params {
		uri = strings.Replace(uri, p.Value, ":"+p.Key, 1)
	}
	return uri
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request) float64 {
	var s int

	s += int(r.ContentLength)

	s += len(r.Proto)
	s += len(r.Method)
	s += len(r.Host)
	s += len(r.RequestURI)

	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}

	return float64(s)
}

func metricReqCounterVec() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method", "handler", "host", "uri"},
	)
}

func metricReqDurHistogramVec() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"code", "method", "uri"},
	)
}

func metricReqSizeSummary() prometheus.Summary {
	return prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystem,
			Name:      "request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		},
	)
}

func metricRespSizeSummary() prometheus.Summary {
	return prometheus.NewSummary(
		prometheus.SummaryOpts{
			Subsystem: subsystem,
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
	)
}
