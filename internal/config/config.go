package config

type Config struct {
	ServerPort string
	RedisAddr  string
	RedisDB    int
}

func Load() *Config {
	return &Config{
		ServerPort: ":8080",
		RedisAddr:  "localhost:6379",
		RedisDB:    0,
	}
}
