package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	router, err := makeRouter(config)
	if err != nil {
		log.Fatalf("could not run server: %v", err)
	}

	_ = http.ListenAndServe(":8080", router)
}

func makeRouter(config *Config) (*mux.Router, error) {
	redisClient := NewRedisClient(config.RedisURL)
	rateLimiter := NewRateLimiter(config.IPLimit, config.TokenLimit, config.ExpirationTimeInSeconds, redisClient)

	router := mux.NewRouter()
	router.Use(RateLimitMiddleware(rateLimiter))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	return router, nil
}
