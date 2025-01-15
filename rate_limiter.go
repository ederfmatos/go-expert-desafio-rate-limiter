package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"net/http"
	"strings"
	"time"
)

type RateLimiter struct {
	ipLimit        int
	tokenLimit     int
	expirationTime time.Duration
	redisClient    *RedisClient
}

func NewRateLimiter(ipLimit, tokenLimit int, expirationTime time.Duration, redisClient *RedisClient) *RateLimiter {
	return &RateLimiter{
		ipLimit:        ipLimit,
		tokenLimit:     tokenLimit,
		expirationTime: expirationTime,
		redisClient:    redisClient,
	}
}

func (rl *RateLimiter) CheckLimit(ctx context.Context, r *http.Request) (bool, error) {
	var ip string

	forwardedFor := r.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		ip = parts[0]
	} else {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}

	token := r.Header.Get("API_KEY")

	var (
		key   string
		limit int
	)

	if token != "" {
		key = token
		limit = rl.tokenLimit
	} else {
		key = ip
		limit = rl.ipLimit
	}

	count, err := rl.redisClient.GetRequestCount(ctx, key)
	if err != nil {
		return false, err
	}

	if count >= limit {
		return false, errors.New("you have reached the maximum number of requests or actions allowed within a certain time frame")
	}

	if err = rl.redisClient.IncrementRequestCount(ctx, key, rl.expirationTime); err != nil {
		return false, fmt.Errorf("increment request count: %v", err)
	}

	return true, err
}
