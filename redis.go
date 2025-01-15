package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(address string) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{Addr: address}),
	}
}

func (c *RedisClient) getRateLimitKey(key string) string {
	return "rate_limit:" + key
}

func (c *RedisClient) IncrementRequestCount(ctx context.Context, key string, expiration time.Duration) error {
	rateLimitKey := c.getRateLimitKey(key)
	if err := c.Client.Incr(ctx, rateLimitKey).Err(); err != nil {
		return fmt.Errorf("increment request count: %v", err)
	}
	c.Client.Expire(ctx, rateLimitKey, expiration)
	return nil
}

func (c *RedisClient) GetRequestCount(ctx context.Context, key string) (int, error) {
	limitKey := c.getRateLimitKey(key)
	count, err := c.Client.Get(ctx, limitKey).Int()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("get request count: %v", err)
	}
	return count, nil
}
