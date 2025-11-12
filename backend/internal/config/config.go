package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Token    string `.env:"TELEGRAM_BOT_TOKEN"`
	Postgres struct {
		Host     string `.env:"PG_HOST"`
		Port     string `.env:"PG_PORT"`
		User     string `.env:"PG_USER"`
		Password string `.env:"PG_PASSWORD"`
		DBName   string `.env:"PG_NAME"`
	}
	Redis struct {
		Host     string        `.env:"RD_HOST"`
		Port     string        `.env:"RD_PORT"`
		User     string        `.env:"RD_USER"`
		Password string        `.env:"RD_PASSWORD"`
		DBName   string        `.env:"RD_NAME"`
		TTL      time.Duration `.env:"RD_TTL"`
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	// Заполняем основную структуру
	config.Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	// Заполняем Postgres
	config.Postgres.Host = os.Getenv("PG_HOST")
	config.Postgres.Port = os.Getenv("PG_PORT")
	config.Postgres.User = os.Getenv("PG_USER")
	config.Postgres.Password = os.Getenv("PG_PASSWORD")
	config.Postgres.DBName = os.Getenv("PG_NAME")

	// Заполняем Redis
	config.Redis.Host = os.Getenv("RD_HOST")
	config.Redis.Port = os.Getenv("RD_PORT")
	config.Redis.User = os.Getenv("RD_USER")
	config.Redis.Password = os.Getenv("RD_PASSWORD")
	config.Redis.DBName = os.Getenv("RD_NAME")

	// Парсим TTL для Redis
	ttlStr := os.Getenv("RD_TTL")
	if ttlStr != "" {
		ttl, err := time.ParseDuration(ttlStr)
		if err != nil {
			return nil, err
		}
		config.Redis.TTL = ttl
	}

	return config, nil
}
