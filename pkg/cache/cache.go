package cache

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheInterface interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Keys(context.Context, string) *redis.StringSliceCmd
	Del(context.Context, ...string) *redis.IntCmd
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":6379",
		Password: "",
		DB:       0,
	})
}
