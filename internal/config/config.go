package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	DB         DBConfig
	Cloudinary CloudinaryConfig
	Postmark   PostmarkConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load env file failed: %w", err)
	}

	cfg := &Config{
		Server:     LoadServerConfig(),
		DB:         LoadDatabaseConfig(),
		Cloudinary: LoadCloudinaryConfig(),
		Postmark:   LoadPostmarkConfig(),
	}

	return cfg, nil
}
