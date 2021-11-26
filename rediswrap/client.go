package rediswrap

import (
	"context"
	"fmt"
	"strings"

	"github.com/DataWorkbench/common/gtrace"
	"github.com/go-redis/redis/v8"
)

const (
	StandaloneMode = "standalone"
	SentinelMode   = "sentinel"
	ClusterMode    = "cluster"
)

type Client interface {
	redis.Cmdable
	AddHook(hook redis.Hook)
	Close() error
}

type RedisConfig struct {
	// Optional Value: "standalone/sentinel/cluster".
	Mode       string `json:"mode"                yaml:"mode"            env:"MODE"            validate:"required"`
	MasterName string `json:"master_name"         yaml:"master_name"     env:"MASTER_NAME"`
	// eg: "127.0.0.1:6379".
	StandaloneAddr string `json:"standalone_addr" yaml:"standalone_addr" env:"STANDALONE_ADDR"`
	// eg: "127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381"
	ClusterAddr string `json:"cluster_addr"       yaml:"cluster_addr"    env:"CLUSTER_ADDR"`
	// eg: "127.0.0.1:7000,127.0.0.1:7001,127.0.0.1:7002,127.0.0.1:7003,127.0.0.1:7004,127.0.0.1:7005".
	SentinelAddr string `json:"sentinel_addr"     yaml:"sentinel_addr"   env:"SENTINEL_ADDR"`
	UserName     string `json:"user_name"         yaml:"user_name"       env:"UER_NAME"`
	Password     string `json:"password"          yaml:"password"        env:"PASSWORD"`
	Database     int    `json:"database"          yaml:"database"        env:"DATABASE"`
}

func NewRedisConn(ctx context.Context, cfg *RedisConfig) (Client, error) {
	var rdb Client
	switch cfg.Mode {
	case StandaloneMode:
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.StandaloneAddr,
			Username: cfg.UserName,
			Password: cfg.Password,
		})
	case SentinelMode:
		rdb = redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			Username:      cfg.UserName,
			SentinelAddrs: strings.Split(cfg.SentinelAddr, ","),
			Password:      cfg.Password,
			DB:            cfg.Database,
		})
	case ClusterMode:
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(cfg.ClusterAddr, ","),
			Username: cfg.UserName,
			Password: cfg.Password,
		})
	default:
		return nil, fmt.Errorf("unsupported mode: %s", cfg.Mode)
	}

	rdb.AddHook(&hookTrace{tracer: gtrace.TracerFromContext(ctx)})
	return rdb, nil
}
