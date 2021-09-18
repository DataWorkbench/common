package rediswrap

import (
	"context"
	"strings"

	"github.com/go-redis/redis/v8"

	"github.com/DataWorkbench/common/gtrace"
)

const (
	SentinelMode   = "sentinel"
	ClusterMode    = "cluster"
	StandaloneMode = "standalone"
)

type Client interface {
	redis.Cmdable
	AddHook(hook redis.Hook)
	Close() error
}

type RedisConfig struct {
	MasterName    string `json:"master_name"         yaml:"master_name"         env:"MASTER_NAME"                    `
	Addr          string `json:"addr"         yaml:"addr"         env:"ADDR"                    `
	Addrs         string `json:"addrs"         yaml:"addrs"         env:"ADDRS"                    `
	UserName      string `json:"user_name"         yaml:"user_name"         env:"UER_NAME"   `
	SentinelAddrs string `json:"sentinel_addrs"         yaml:"sentinel_addrs"         env:"SENTINEL_ADDRS"                    `
	Password      string `json:"password"      yaml:"password"      env:"PASSWORD"                  `
	Database      int    `json:"database"      yaml:"database"      env:"DATABASE"                  `
	Mode          string `json:"mode"      yaml:"mode"      env:"MODE"                  `
}

func NewRedisConn(ctx context.Context, cfg *RedisConfig) (Client, error) {
	var rdb Client
	switch cfg.Mode {
	case SentinelMode:
		rdb = redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:    cfg.MasterName,
			Username:      cfg.UserName,
			SentinelAddrs: strings.Split(cfg.SentinelAddrs, ","),
			Password:      cfg.Password,
			DB:            cfg.Database,
		})
	case StandaloneMode:
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Username: cfg.UserName,
			Password: cfg.Password,
		})
	case ClusterMode:
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    strings.Split(cfg.Addrs, ","),
			Username: cfg.UserName,
			Password: cfg.Password,
		})
	default:
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.Addr,
			Username: cfg.UserName,
			Password: cfg.Password,
		})
	}
	rdb.AddHook(&hookTrace{tracer: gtrace.TracerFromContext(ctx)})
	return rdb, nil
}
