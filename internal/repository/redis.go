package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/ttymayor/url-shortener/internal/config"
)

var ctx = context.Background()

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// health check / ping
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return rdb
}
