package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLimiter struct {
	client *redis.Client
	rate   int
	window time.Duration
}

func NewRedisLimiter(addr, password string, db int, rate int, window time.Duration) (Limiter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &redisLimiter{client: client, rate: rate, window: window}, nil
}

func NewRedisLimiterFromClient(client *redis.Client, rate int, window time.Duration) Limiter {
	return &redisLimiter{client: client, rate: rate, window: window}
}

func (rl *redisLimiter) Allow(key string) bool {
	ctx := context.Background()
	fullKey := "rate:" + key
	ttl := int(rl.window.Seconds())

	count, err := rl.client.Incr(ctx, fullKey).Result()
	if err != nil {
		return true
	}
	if count == 1 {
		rl.client.Expire(ctx, fullKey, time.Duration(ttl)*time.Second)
	}
	return int(count) <= rl.rate
}
