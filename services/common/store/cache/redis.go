package cache

import (
	"context"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

type Redis struct {
	client *redis.Client
}

func New(ctx context.Context) *Redis {
	zap.L().Info("Creating new Redis client")

	client := redis.NewClient(&redis.Options{
		// TODO: Add additional configuration here
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}).WithContext(ctx)

	_, err := client.Ping().Result()

	if err != nil {
		zap.L().Sugar().Fatalf("Failed to connect to Redis", err)
	}

	zap.L().Info("Connected to Redis")
	return &Redis{
		client: client,
	}
}

func (r *Redis) Close() {
	r.client.Close()
}

func (r *Redis) GetClient(ctx context.Context) *redis.Client {
	return r.client.WithContext(ctx)
}
