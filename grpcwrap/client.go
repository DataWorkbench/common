package grpcwrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/DataWorkbench/glog"
	"github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/keepalive"
)

// ClientConfig used to create an connection to grpc server
type ClientConfig struct {
	// Address sample "127.0.0.1:50001" or "127.0.0.1:50001, 127.0.0.1:50002, 127.0.0.1:50003"
	Address string `json:"address" yaml:"address" env:"ADDRESS" validate:"required"`
	// grpc log level, 1 => info, 2 => waring, 3 => error, 4 => fatal
	LogLevel     int `json:"log_level"     yaml:"log_level"     env:"LOG_LEVEL,default=2"     validate:"gte=1,lte=4"`
	LogVerbosity int `json:"log_verbosity" yaml:"log_verbosity" env:"LOG_VERBOSITY,default=1" validate:"required"`
}

// NewConn return an new grpc.ClientConn
// NOTICE: Must set glog.Logger into the ctx by glow.WithContext
func NewConn(ctx context.Context, cfg *ClientConfig) (conn *grpc.ClientConn, err error) {
	lp := glog.FromContext(ctx)

	defer func() {
		if err != nil {
			lp.Error().Error("connected to grpc server error", err).Fire()
		}
	}()

	lp.Info().Msg("connecting to grpc server").String("address", cfg.Address).Fire()

	// TODO: support balance
	// address format "127.0.0.1:50001" or "127.0.0.1:50001, 127.0.0.1:50002, 127.0.0.1:50003"
	hosts := strings.Split(strings.ReplaceAll(cfg.Address, " ", ""), ",")
	if len(hosts) == 0 {
		err = fmt.Errorf("invalid address: %s", cfg.Address)
		return
	}

	// setup grpc logger
	grpclog.SetLoggerV2(&Logger{
		Output:    lp,
		Verbosity: cfg.LogVerbosity,
		Level:     cfg.LogLevel,
	})

	var opts []grpc.DialOption
	// Set and add insecure
	opts = append(opts, grpc.WithInsecure())

	// set and add connect params
	opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff: backoff.Config{
			BaseDelay:  time.Millisecond * 100, // Default was 1s.
			Multiplier: 1.6,                    // Default
			Jitter:     0.2,                    // Default
			MaxDelay:   time.Second * 3,        // Default was 120s.
		},
		MinConnectTimeout: time.Second * 5,
	}))

	// Setup keepalive params
	opts = append(opts, grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                time.Second * 30,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		},
	))

	// Set and add Unary Client Interceptor
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(0),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second*1)),
		grpc_retry.WithCodes(codes.Unavailable, codes.Aborted, codes.DeadlineExceeded, codes.ResourceExhausted),
	)))

	opts = append(opts, grpc.WithChainUnaryInterceptor(
		grpc_prometheus.UnaryClientInterceptor,
		basicUnaryClientInterceptor(),
	))

	// TODO: Impls and add Stream Client Interceptor

	conn, err = grpc.DialContext(ctx, hosts[0], opts...)
	return
}
