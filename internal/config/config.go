package config

type Config struct {
	ServerPort string
}

func Load() *Config {
	return &Config{
		ServerPort: ":8080",
	}
}
