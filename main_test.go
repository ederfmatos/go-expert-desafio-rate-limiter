package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterWithRedis(t *testing.T) {
	ctx := context.Background()

	redisContainer, err := redis.Run(ctx, "docker.io/redis:7")
	require.NoError(t, err)
	defer redisContainer.Terminate(ctx)

	redisPort, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	config := &Config{
		IPLimit:                 5,
		TokenLimit:              10,
		ExpirationTimeInSeconds: 10,
		RedisURL:                "localhost:" + redisPort.Port(),
	}
	router, err := makeRouter(config)
	require.NoError(t, err)

	server := httptest.NewServer(router)
	defer server.Close()

	t.Run("Rate Limit with Token", func(t *testing.T) {
		for i := 0; i <= config.TokenLimit; i++ {
			resp := makeRequest(t, server.URL, "abc123", "")

			if i < config.TokenLimit {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				checkResponseBody(t, resp, "OK")
			} else {
				require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}
	})

	t.Run("Rate Limit with IP", func(t *testing.T) {
		for i := 0; i <= config.IPLimit; i++ {
			resp := makeRequest(t, server.URL, "", "192.168.1.1:8080")

			if i < config.IPLimit {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				checkResponseBody(t, resp, "OK")
			} else {
				require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}
	})

	t.Run("Rate Limit with IP waiting expiration", func(t *testing.T) {
		remoteAddress := "192.168.1.2:8080"

		for i := 0; i < config.IPLimit; i++ {
			resp := makeRequest(t, server.URL, "", remoteAddress)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			checkResponseBody(t, resp, "OK")
		}

		resp := makeRequest(t, server.URL, "", remoteAddress)
		require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

		time.Sleep(time.Duration(config.TokenLimit) * time.Second)

		resp = makeRequest(t, server.URL, "", remoteAddress)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		checkResponseBody(t, resp, "OK")
	})

	t.Run("Rate Limit with IP and Token", func(t *testing.T) {
		remoteAddress := "192.168.1.3:8080"

		for i := 0; i <= config.TokenLimit; i++ {
			resp := makeRequest(t, server.URL, "abc123", remoteAddress)

			if i < config.TokenLimit {
				require.Equal(t, http.StatusOK, resp.StatusCode)
				checkResponseBody(t, resp, "OK")
			} else {
				require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
			}
		}
	})

	t.Run("Rate Limit with Different IP", func(t *testing.T) {
		remoteAddress := "192.168.1.4:8080"

		for i := 0; i < config.IPLimit; i++ {
			resp := makeRequest(t, server.URL, "", remoteAddress)
			require.Equal(t, http.StatusOK, resp.StatusCode)
			checkResponseBody(t, resp, "OK")
		}

		resp := makeRequest(t, server.URL, "", remoteAddress)
		require.Equal(t, http.StatusTooManyRequests, resp.StatusCode)
	})
}

func makeRequest(t *testing.T, url string, apiKey, remoteAddress string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	if apiKey != "" {
		req.Header.Set("API_KEY", apiKey)
	}
	if remoteAddress != "" {
		req.Header.Set("X-Forwarded-For", remoteAddress)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func checkResponseBody(t *testing.T, resp *http.Response, expectedBody string) {
	responseBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, expectedBody, string(responseBody))
}
