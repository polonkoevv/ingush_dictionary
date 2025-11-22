package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Token   string
	Logging struct {
		LogFile  string
		LogLevel string
	}
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
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

	// Заполняем Logging
	config.Logging.LogFile = os.Getenv("LOG_FILE")
	if config.Logging.LogFile == "" {
		config.Logging.LogFile = "logs/app.log" // значение по умолчанию
	}
	config.Logging.LogLevel = os.Getenv("LOG_LEVEL")
	if config.Logging.LogLevel == "" {
		config.Logging.LogLevel = "info" // значение по умолчанию
	}

	// Заполняем Postgres
	config.Postgres.Host = os.Getenv("PG_HOST")
	config.Postgres.Port = os.Getenv("PG_PORT")
	config.Postgres.User = os.Getenv("PG_USER")
	config.Postgres.Password = os.Getenv("PG_PASSWORD")
	config.Postgres.DBName = os.Getenv("PG_NAME")
	return config, nil
}
