package getcd

import (
	"context"
	"strings"
	"time"

	"github.com/DataWorkbench/glog"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	etcdv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Client = etcdv3.Client

// Config is a copy of clientv3.Config.
type Config struct {
	// Etcd server endpoints, multiple endpoint are separated by ",".
	// eg: "127.0.0.1:2379" or "127.0.0.1:2379,127.0.0.2:2379,127.0.0.3:2379".
	Endpoints   string        `json:"endpoints" yaml:"endpoints" env:"ENDPOINTS" validate:"required"`
	DialTimeout time.Duration `json:"dial_timeout" yaml:"dial_timeout" env:"DIAL_TIMEOUT,default=5s" validate:"required"`
}

// NewClient creates a new etcd Client.
func NewClient(ctx context.Context, cfg *Config, options ...Option) (cli *Client, err error) {
	opts := applyOptions(options...)
	lg := glog.FromContext(ctx)

	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithChainUnaryInterceptor(
		otgrpc.OpenTracingClientInterceptor(opts.tracer),
		grpc_prometheus.UnaryClientInterceptor,
	))
	dialOpts = append(dialOpts, grpc.WithChainStreamInterceptor(
		otgrpc.OpenTracingStreamClientInterceptor(opts.tracer),
		grpc_prometheus.StreamClientInterceptor,
	))

	lg.Debug().String("connecting to etcd endpoints", cfg.Endpoints).Fire()
	cli, err = etcdv3.New(etcdv3.Config{
		Endpoints:   strings.Split(cfg.Endpoints, ","),
		DialTimeout: cfg.DialTimeout,
		DialOptions: dialOpts,
	})
	if err != nil {
		lg.Error().Error("connects to etcd error", err).Fire()
		return
	}
	lg.Debug().Msg("successful connection to etcd").Fire()
	return
}
