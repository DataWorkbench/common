package rediswrap

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	// Hosts sample "127.0.0.1:6379"
	Hosts    string `json:"hosts"         yaml:"hosts"         env:"HOSTS"                     validate:"required"`
	Password string `json:"password"      yaml:"password"      env:"PASSWORD"                  `
	Database int    `json:"database"      yaml:"database"      env:"DATABASE"                  `
}

func NewRedisConn(ctx context.Context, cfg *RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Hosts,
		Password: cfg.Password,
		DB:       cfg.Database,
	})
	return rdb, nil

}
