package main

import (
	"context"
	"log"
	"test/internal/application"
	"test/internal/config"
	"test/internal/handler/tg"
	"test/internal/infrustructure/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	ctx := context.Background()

	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewPostgresConnection(cfg)

	if err != nil {
		log.Fatal(err)
	}

	repository, err := postgres.CreateRepositories(db)

	if err != nil {
		log.Fatal(err)
	}

	service, err := application.CreateServices(repository)

	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)

	if err != nil {
		log.Panic(err)
	}

	tgHandler := tg.NewTgHandler(bot, service)

	bot.Debug = true
	// Передаём управление универсальному обработчику
	tgHandler.Run(ctx)
}
