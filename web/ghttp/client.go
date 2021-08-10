package ghttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/opentracing/opentracing-go"

	"github.com/DataWorkbench/common/gtrace"
)

// ClientConfig for configuration http transport
type ClientConfig struct {
	ExpectContinueTimeout time.Duration `json:"expect_continue_timeout" yaml:"expect_continue_timeout" env:"EXPECT_CONTINUE_TIMEOUT,default=2s"  validate:"required"`
	ResponseHeaderTimeout time.Duration `json:"response_header_timeout" yaml:"response_header_timeout" env:"RESPONSE_HEADER_TIMEOUT,default=0s"  validate:"required"`
	TLSHandshakeTimeout   time.Duration `json:"tls_handshake_timeout"   yaml:"tls_handshake_timeout"   env:"TLS_HANDSHAKE_TIMEOUT,default=10s"   validate:"required"`
	IdleConnTimeout       time.Duration `json:"idle_conn_timeout"       yaml:"idle_conn_timeout"       env:"IDLE_CONN_TIMEOUT,default=30s"       validate:"required"`
	MaxIdleConns          int           `json:"max_idle_conns"          yaml:"max_idle_conns"          env:"MAX_IDLE_CONNS,default=128"          validate:"required"`
	MaxIdleConnsPerHost   int           `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host" env:"MAX_IDLE_CONNS_PER_HOST,default=128" validate:"required"`
}

// defaultClientConfig return a ClientConfig that can be used in most scenarios;
// And without TLSClientConfig and Proxy
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		ExpectContinueTimeout: time.Second * 2,
		ResponseHeaderTimeout: 0, // time.Second * 15
		TLSHandshakeTimeout:   time.Second * 10,
		IdleConnTimeout:       time.Second * 30,
		MaxIdleConns:          128,
		MaxIdleConnsPerHost:   128,
	}
}

type Client struct {
	*http.Client
	tracer opentracing.Tracer
}

// NewClient creating a new http.Client using the provided NetDialer
func NewClient(ctx context.Context, cfg *ClientConfig) *Client {
	if cfg == nil {
		cfg = NewClientConfig()
	}
	dialer := &net.Dialer{
		Timeout:       time.Second * 10,
		KeepAlive:     time.Second * 15,
		DualStack:     true,
		LocalAddr:     nil,
		FallbackDelay: 0,
		Resolver:      nil,
		Cancel:        nil,
		Control:       nil,
	}

	transport := &http.Transport{
		Proxy:                 nil,
		TLSClientConfig:       nil,
		DialContext:           dialer.DialContext,
		ExpectContinueTimeout: cfg.ExpectContinueTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
		TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
	}

	cli := &Client{
		Client: &http.Client{
			Transport: transport,
			Timeout:   0,
		},
		tracer: gtrace.TracerFromContext(ctx),
	}
	return cli
}

// Send is wrapper for http.Client.Do. To support opentracing span.
func (cli *Client) Send(ctx context.Context, req *http.Request) (resp *http.Response, err error) {
	lg := glog.FromContext(ctx)

	if tid := gtrace.IdFromContext(ctx); tid != "" {
		req.Header.Set(gtrace.HeaderKey, tid)
	}
	if span := opentracing.SpanFromContext(ctx); span != nil {
		err = cli.tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		if err != nil {
			lg.Error().Error("inject span to request header error", err).Fire()
		}
	}

	lg.Debug().String("sending request to url", req.URL.String()).Fire()
	resp, err = cli.Client.Do(req)
	if err != nil {
		lg.Error().Error("send request error", err).Fire()
		return
	}
	lg.Debug().Int("successful request with status", resp.StatusCode).Fire()
	return
}
