package main

import (
	"log"
	"test/internal/handler/http"
	"test/internal/handler/tg"
	"test/internal/handler/tg/utils"
	"test/internal/service"
	"test/internal/storage/inmemory"
	_ "test/internal/storage/inmemory"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	bot, err := tgbotapi.NewBotAPI("8140364030:AAEuULwCMW1MUu9vef7CZpfRuhHqh-QHeRo")
	if err != nil {
		log.Panic(err)
	}

	storage := inmemory.NewInMemoryStorage()
	service := service.NewLanguageService(storage)
	tgHandler := tg.NewTgHandler(bot, storage)
	httpHandler := http.NewHttpHandler(service)

	bot.Debug = true

	// log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Выбрать язык"),
		),
	)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		// if update.Message.Audio != nil || update.Message.Video != nil || update.Message.Voice != nil || update.Message.Document != nil || update.Message.VideoNote != nil {
		if update.Message.Text == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный тип сообщения")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			continue
		}
		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для перевода слов с русского на ингушский и обратно. Чтобы начать, нажмите кнопку 'Сменить язык'.")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			continue
		}

		if update.Message.Text == "Выбрать язык" || update.Message.Text == "Сменить язык" || update.Message.Text == "RUS -> ING" || update.Message.Text == "ING -> RUS" {
			language, err := tgHandler.ChangeLanguage(update.Message.Chat.ID)
			if err != nil {
				log.Panic(err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбранный язык: русский")
			keyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("RUS -> ING"),
				),
			)
			if language == "ing" {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбранный язык: ингушский")
				keyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("ING -> RUS"),
					),
				)
			}

			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		} else if update.Message.Text == "Audio" {
			voiceMsg := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FilePath("C:\\Users\\Admin\\Desktop\\Coding\\testbot\\5fbd70e24033a.ogg"))
			// voiceMsg.Caption = "test"
			// voiceMsg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			// 	tgbotapi.NewInlineKeyboardRow(
			// 		tgbotapi.NewInlineKeyboardButtonData("test", "/start"),
			// 	),
			// )
			// voiceMsg := tgbotapi.NewVoice(update.Message.Chat.ID, tgbotapi.FileURL("doshlorg.ru/storage/5fbd70e24033a.ogg"))
			bot.Send(voiceMsg)
		} else if update.Message.Text == "test" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выбранный язык: ингушский")
			keyboard := utils.GetButtons(10, 10, "test", "test")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		} else {
			language, err := service.GetLanguage(update.Message.Chat.ID)
			if err != nil {
				log.Panic(err)
			}

			res, total, err := httpHandler.GetWord(update.Message.Text, language)
			if err != nil {
				log.Panic(err)
			}
			if res == "" {
				if language == "rus" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Слово не найдено")
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Цу тайпара дош дац")
					msg.ReplyToMessageID = update.Message.MessageID
					msg.ReplyMarkup = keyboard
					bot.Send(msg)
				}
				continue
			}
			log.Println(res)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
			msg.ReplyToMessageID = update.Message.MessageID
			keyboard := utils.GetButtons(total, total, language, update.Message.Text)
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}
	}
}
