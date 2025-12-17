package config

type FailureStrategy string

const (
	FailOpen   FailureStrategy = "fail_open"
	FailClosed FailureStrategy = "fail_closed"
)

type Config struct {
	ServerPort      string
	RedisAddr       string
	RedisDB         int
	FailureStrategy FailureStrategy
}

func Load() *Config {
	return &Config{
		ServerPort:      ":8080",
		RedisAddr:       "localhost:6379",
		RedisDB:         0,
		FailureStrategy: FailClosed, // default for public APIs
	}
}
