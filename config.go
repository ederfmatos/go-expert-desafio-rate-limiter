package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

type Config struct {
	IPLimit                 int
	TokenLimit              int
	ExpirationTimeInSeconds time.Duration
	RedisURL                string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load environment variables: %v", err)
	}

	ipLimit, err := getIntEnv("IP_LIMIT")
	if err != nil {
		return nil, fmt.Errorf("parse IP_LIMIT: %v", err)
	}

	tokenLimit, err := getIntEnv("TOKEN_LIMIT")
	if err != nil {
		return nil, fmt.Errorf("parse TOKEN_LIMIT: %v", err)
	}

	expirationTime, err := getIntEnv("EXPIRATION_TIME_IN_SECONDS")
	if err != nil {
		return nil, fmt.Errorf("parse EXPIRATION_TIME_IN_SECONDS: %v", err)
	}

	return &Config{
		IPLimit:                 ipLimit,
		TokenLimit:              tokenLimit,
		ExpirationTimeInSeconds: time.Duration(expirationTime) * time.Second,
		RedisURL:                os.Getenv("REDIS_URL"),
	}, nil
}

func getIntEnv(name string) (int, error) {
	return strconv.Atoi(os.Getenv(name))
}
