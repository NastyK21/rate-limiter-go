package limiter

import (
	"context"
	"os"
)

type RateLimiter struct {
	redis  *RedisClient
	script string
}

func NewRateLimiter(redis *RedisClient) (*RateLimiter, error) {
	data, err := os.ReadFile("internal/limiter/lua/token_bucket.lua")
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		redis:  redis,
		script: string(data),
	}, nil
}

func (r *RateLimiter) Allow(ctx context.Context, key string, capacity, refillRate float64) (bool, float64, error) {
	now, err := r.redis.Client.Time(ctx).Result()
	if err != nil {
		return false, 0, err
	}

	result, err := r.redis.Client.Eval(
		ctx,
		r.script,
		[]string{key},
		capacity,
		refillRate,
		now.Unix(),
	).Result()

	if err != nil {
		return false, 0, err
	}

	res := result.([]interface{})
	allowed := res[0].(int64) == 1
	var remaining float64

	switch v := res[1].(type) {
	case int64:
		remaining = float64(v)
	case float64:
		remaining = v
	default:
		remaining = 0
	}

	return allowed, remaining, nil
}
