package tg

import (
	"context"
	"log"
	"strings"
	"test/internal/application"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgHandler struct {
	bot *tgbotapi.BotAPI
	srv *application.Service
}

func NewTgHandler(bot *tgbotapi.BotAPI, srv *application.Service) *TgHandler {
	return &TgHandler{bot: bot, srv: srv}
}

func setupBotCommands(bot *tgbotapi.BotAPI) error {
	// Настройка списка команд
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Перезапустить бота"},
		{Command: "language", Description: "Настройка языка"},
		{Command: "dictionaries", Description: "Список словарей"},
		{Command: "choose", Description: "Выбор словарей"},
		{Command: "help", Description: "Помощь и команды"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.Request(config)
	return err
}

// Run запускает цикл обработки апдейтов TG, включая сообщения и callback-кнопки.
func (h *TgHandler) Run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	if err := setupBotCommands(h.bot); err != nil {
		log.Printf("Ошибка настройки команд: %v", err)
	}

	updates := h.bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil {
			switch {
			case strings.HasPrefix(update.CallbackQuery.Data, "set_dict_"):
				err := h.addDictCallback(ctx, update.CallbackQuery)

				if err != nil {
					h.errorCallback(ctx, update.CallbackQuery)
				}
				continue
			case strings.HasPrefix(update.CallbackQuery.Data, "del_dict_"):
				err := h.remDictCallback(ctx, update.CallbackQuery)
				if err != nil {
					h.errorCallback(ctx, update.CallbackQuery)
				}
				continue
			case strings.HasPrefix(update.CallbackQuery.Data, "del_msg"):
				err := h.deleteMessageCallback(ctx, update.CallbackQuery)
				if err != nil {
					h.errorCallback(ctx, update.CallbackQuery)
				}
				continue
			}
		}

		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		switch update.Message.Text {
		case "/start":
			h.hello(ctx, &update)
			h.help(ctx, &update)
			h.deleteMessageSafe(update.Message.Chat.ID, update.Message.MessageID)
			continue
		case "/help":
			h.help(ctx, &update)
			h.deleteMessageSafe(update.Message.Chat.ID, update.Message.MessageID)
			continue
		case "/language":
			h.deleteMessageSafe(update.Message.Chat.ID, update.Message.MessageID)
			h.changeLanguage(ctx, &update)
			continue
		case "/dictionaries":
			h.deleteMessageSafe(update.Message.Chat.ID, update.Message.MessageID)
			h.listDictionaries(ctx, &update)
			continue
		case "/choose":
			h.deleteMessageSafe(update.Message.Chat.ID, update.Message.MessageID)
			h.chooseDicts(ctx, &update)
			continue
		default:
			h.translate(ctx, &update)
			continue

		}
	}
}
