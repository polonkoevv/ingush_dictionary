package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test/internal/application"
	"test/internal/config"
	"test/internal/handler/tg"
	"test/internal/infrastructure/postgres"
	"test/internal/logger"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Error while loading config", slog.Any("error", err))
		os.Exit(1)
	}

	// Настраиваем логгер с записью в файл
	appLogger, closeLogFile, err := logger.SetupLogger(cfg.Logging.LogFile, cfg.Logging.LogLevel)
	if err != nil {
		slog.Error("Error while setting up logger", slog.Any("error", err))
		os.Exit(1)
	}
	slog.SetDefault(appLogger)
	defer closeLogFile()

	slog.Info("Loaded config")

	db, err := postgres.NewPostgresConnection(cfg)
	if err != nil {
		slog.Error("Error while creating DB connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error closing database connection", slog.Any("error", err))
		} else {
			slog.Info("Database connection closed")
		}
	}()
	slog.Info("Connected to database")

	repository, err := postgres.CreateRepositories(db)
	if err != nil {
		slog.Error("Error while creating repositories", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Created repositories")

	service, err := application.CreateServices(repository)
	if err != nil {
		slog.Error("Error while creating services", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Created services")

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		slog.Error("Error while creating Telegram Bot", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("Created Telegram bot")

	tgHandler := tg.NewTgHandler(bot, service)

	// Создаем context с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем обработку сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Запускаем обработчик в горутине
	done := make(chan error, 1)
	go func() {
		slog.Info("System is ready...")
		done <- tgHandler.Run(ctx)
	}()

	// Ждем сигнал завершения или ошибку
	select {
	case sig := <-sigChan:
		slog.Info("Received shutdown signal", slog.String("signal", sig.String()))
		cancel() // Отменяем context

		// Даем время на завершение работы (graceful shutdown timeout)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Ждем завершения или таймаут
		select {
		case err := <-done:
			if err != nil {
				slog.Error("Handler exited with error", slog.Any("error", err))
			} else {
				slog.Info("Handler stopped gracefully")
			}
		case <-shutdownCtx.Done():
			slog.Warn("Shutdown timeout reached, forcing exit")
		}
	case err := <-done:
		if err != nil {
			slog.Error("Handler exited with error", slog.Any("error", err))
			os.Exit(1)
		}
		slog.Info("Handler stopped")
	}

	slog.Info("Application shutdown complete")
}
