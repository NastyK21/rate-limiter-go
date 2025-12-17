package limiter

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(_ string, _ int) (*RedisClient, error) {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return nil, redis.ErrClosed
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	opts.DialTimeout = 5 * time.Second
	opts.ReadTimeout = 3 * time.Second
	opts.WriteTimeout = 3 * time.Second
	opts.PoolSize = 10
	opts.MinIdleConns = 2

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisClient{Client: rdb}, nil
}
