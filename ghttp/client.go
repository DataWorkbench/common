package ghttp

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

// ClientConfig for configuration http transport
type ClientConfig struct {
	ExpectContinueTimeout time.Duration `json:"expect_continue_timeout" yaml:"expect_continue_timeout" env:"EXPECT_CONTINUE_TIMEOUT,default=2s"  validate:"required"`
	ResponseHeaderTimeout time.Duration `json:"response_header_timeout" yaml:"response_header_timeout" env:"RESPONSE_HEADER_TIMEOUT,default=0s"  validate:"required"`
	TLSHandshakeTimeout   time.Duration `json:"tls_handshake_timeout"   yaml:"tls_handshake_timeout"   env:"TLS_HANDSHAKE_TIMEOUT,default=10s"   validate:"required"`
	IdleConnTimeout       time.Duration `json:"idle_conn_timeout"       yaml:"idle_conn_timeout"       env:"IDLE_CONN_TIMEOUT,default=30s"       validate:"required"`
	MaxIdleConns          int           `json:"max_idle_conns"          yaml:"max_idle_conns"          env:"MAX_IDLE_CONNS,default=128"          validate:"required"`
	MaxIdleConnsPerHost   int           `json:"max_idle_conns_per_host" yaml:"max_idle_conns_per_host" env:"MAX_IDLE_CONNS_PER_HOST,default=128" validate:"required"`
	// TODO:
	TLSClientConfig *tls.Config                           `json:"-" yaml:"-" env:"-" validate:"-"`
	Proxy           func(*http.Request) (*url.URL, error) `json:"-" yaml:"-" env:"-" validate:"-"`
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
		TLSClientConfig:       nil,
		Proxy:                 nil,
	}
}

// NewClient creating a new http.Client using the provided NetDialer
func NewClient(cfg *ClientConfig) *http.Client {
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
		Proxy:                 cfg.Proxy,
		TLSClientConfig:       cfg.TLSClientConfig,
		DialContext:           dialer.DialContext,
		ExpectContinueTimeout: cfg.ExpectContinueTimeout,
		ResponseHeaderTimeout: cfg.ResponseHeaderTimeout,
		TLSHandshakeTimeout:   cfg.TLSHandshakeTimeout,
		IdleConnTimeout:       cfg.IdleConnTimeout,
		MaxIdleConns:          cfg.MaxIdleConns,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnsPerHost,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   0,
	}

	return client
}
