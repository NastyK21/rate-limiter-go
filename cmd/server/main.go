package main

import (
	"log"
	"net/http"

	"github.com/NastyK21/rate-limiter-go/internal/config"
	"github.com/NastyK21/rate-limiter-go/internal/limiter"
	"github.com/NastyK21/rate-limiter-go/internal/metrics"
	"github.com/NastyK21/rate-limiter-go/internal/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	metrics.Register()

	cfg := config.Load()

	redisClient, err := limiter.NewRedisClient(cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Printf("redis unavailable, starting in degraded mode: %v", err)
		redisClient = nil
	}

	rl := limiter.NewRateLimiter(redisClient)
	localLimiter := limiter.NewLocalLimiter(5, 1)

	rateLimitCfg := middleware.RateLimitConfig{
		// IP limits (anonymous)
		Capacity:   5,
		RefillRate: 1,

		// User limits (authenticated)
		UserCapacity:   20,
		UserRefillRate: 5,

		FailureStrategy: cfg.FailureStrategy,
		LocalLimiter:    localLimiter,
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

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
