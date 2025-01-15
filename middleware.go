package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RateLimitMiddleware(rateLimiter *RateLimiter) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok, err := rateLimiter.CheckLimit(r.Context(), r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusTooManyRequests)
				return
			}
			if !ok {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
			handler.ServeHTTP(w, r)
		})
	}
}
