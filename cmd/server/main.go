package main

import (
	"context"
	"log"
	"net/http"

	"github.com/NastyK21/rate-limiter-go/internal/config"
	"github.com/NastyK21/rate-limiter-go/internal/limiter"
)

func main() {
	cfg := config.Load()

	redisClient, err := limiter.NewRedisClient(cfg.RedisAddr, cfg.RedisDB)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	// -------- Phase 3: TEMP limiter test --------
	rl, err := limiter.NewRateLimiter(redisClient)
	if err != nil {
		log.Fatal(err)
	}

	allowed, remaining, err := rl.Allow(
		context.Background(),
		"test-user",
		5, // capacity
		1, // refill rate (tokens/sec)
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("allowed:", allowed, "remaining:", remaining)
	// -------- END TEMP test --------

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: mux,
	}

	log.Println("server started on", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
