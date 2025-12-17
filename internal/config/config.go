package config

import (
	"log"
	"os"
	"strings"
)

type FailureStrategy string

const (
	FailOpen   FailureStrategy = "fail-open"
	FailClosed FailureStrategy = "fail-closed"
)

type Config struct {
	ServerPort      string
	RedisAddr       string
	RedisDB         int
	FailureStrategy FailureStrategy
}

func Load() Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	failure := FailureStrategy(os.Getenv("FAILURE_STRATEGY"))
	if failure != FailOpen && failure != FailClosed {
		failure = FailOpen
	}

	log.Println("CONFIG:")
	log.Println("  SERVER_PORT =", port)
	log.Println("  REDIS_ADDR  =", redisAddr)
	log.Println("  FAILURE_STRATEGY =", failure)

	return Config{
		ServerPort:      port,
		RedisAddr:       strings.TrimSpace(redisAddr),
		RedisDB:         0,
		FailureStrategy: failure,
	}
}
