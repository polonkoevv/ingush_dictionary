package main

import (
	"log"
	"os"
	"test/internal/handler/tg"
	"test/internal/service"
	"test/internal/storage/inmemory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("warning: TELEGRAM_BOT_TOKEN is empty")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	storage := inmemory.NewInMemoryStorage()
	service := service.NewLanguageService(storage)
	tgHandler := tg.NewTgHandler(bot, storage, service)

	bot.Debug = true
	// Передаём управление универсальному обработчику
	tgHandler.Run()
}
