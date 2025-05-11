package db

import (
	"context"

	"fliQt/config"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}

func PingRedis(ctx context.Context, client *redis.Client) error {
	_, err := client.Ping(ctx).Result()
	return err
}
