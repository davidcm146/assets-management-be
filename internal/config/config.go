package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("load env file failed: %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: os.Getenv("PORT"),
		},
		DB: DBConfig{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			DBName:   os.Getenv("DB_NAME"),
		},
	}

	return cfg, nil
}
