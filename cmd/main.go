package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"test/internal/application"
	"test/internal/handler/tg"
	"test/internal/infrustructure/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	db, err := postgres.NewPostgresConnection()

	if err != nil {
		log.Fatal(err)
	}

	usrRep, wordRep, dictRep, err := postgres.CreateRepositories(db)

	usrSrv, wrdSrv, dictSrv, err := application.CreateServices(&usrRep, &wordRep, &dictRep)

	fmt.Println(err)

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Println("warning: TELEGRAM_BOT_TOKEN is empty")
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	tgHandler := tg.NewTgHandler(bot, *usrSrv, *wrdSrv, *dictSrv)

	bot.Debug = true
	// Передаём управление универсальному обработчику
	tgHandler.Run(ctx)
}
