package main

import (
	"log"
	"net/http"

	"github.com/NastyK21/rate-limiter-go/internal/config"
	"github.com/NastyK21/rate-limiter-go/internal/limiter"
	"github.com/NastyK21/rate-limiter-go/internal/middleware"
)

func main() {
	cfg := config.Load()

	redisClient, err := limiter.NewRedisClient(cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	rl, err := limiter.NewRateLimiter(redisClient)
	if err != nil {
		log.Fatal(err)
	}

	rateLimitCfg := middleware.RateLimitConfig{
		Capacity:   5,
		RefillRate: 1,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := middleware.RateLimit(rl, rateLimitCfg)(mux)

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: handler,
	}

	log.Println("server started on", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
