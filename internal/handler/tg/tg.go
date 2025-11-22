package tg

import (
	"context"
	"log/slog"
	"strings"
	"test/internal/application"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgHandler struct {
	bot            *tgbotapi.BotAPI
	srv            *application.Service
	messageCleaner *MessageCleaner
	messageTTL     time.Duration
}

func NewTgHandler(bot *tgbotapi.BotAPI, srv *application.Service, messageTTL time.Duration) *TgHandler {
	cleaner := NewMessageCleaner(bot)
	return &TgHandler{
		bot:            bot,
		srv:            srv,
		messageCleaner: cleaner,
		messageTTL:     messageTTL,
	}
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
func (h *TgHandler) Run(ctx context.Context) error {
	// Запускаем cleaner в фоне
	go h.messageCleaner.Start(ctx)
	defer h.messageCleaner.Stop()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	if err := setupBotCommands(h.bot); err != nil {
		slog.Error("Ошибка настройки команд",
			slog.String("component", "tg_handler"),
			slog.String("op", "setupBotCommands"),
			slog.Any("error", err))
	}

	updates := h.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down Telegram handler",
				slog.String("component", "tg_handler"),
				slog.String("op", "Run"))
			// Останавливаем канал обновлений
			h.bot.StopReceivingUpdates()
			return ctx.Err()
		case update, ok := <-updates:
			if !ok {
				slog.Info("Updates channel closed",
					slog.String("component", "tg_handler"),
					slog.String("op", "Run"))
				return nil
			}

			if update.CallbackQuery != nil {
				switch {
				case strings.HasPrefix(update.CallbackQuery.Data, "set_dict"):
					err := h.setDictCallback(ctx, update.CallbackQuery)

					if err != nil {
						slog.Error("telegram error", slog.String("component", "tg_handler"),
							slog.String("op", "setDictCallback"),
							slog.Any("error", err),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
						h.errorCallback(ctx, update.CallbackQuery)
					} else {
						slog.Info("telegram callback processed", slog.String("component", "tg_handler"),
							slog.String("op", "setDictCallback"),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
					}

					continue
				case strings.HasPrefix(update.CallbackQuery.Data, "rem_dict"):
					err := h.remDictCallback(ctx, update.CallbackQuery)
					if err != nil {

						slog.Error("telegram error", slog.String("component", "tg_handler"),
							slog.String("op", "remDictCallback"),
							slog.Any("error", err),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)

						h.errorCallback(ctx, update.CallbackQuery)
					} else {
						slog.Info("telegram callback processed", slog.String("component", "tg_handler"),
							slog.String("op", "remDictCallback"),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
					}
					continue
				case strings.HasPrefix(update.CallbackQuery.Data, "change_page"):
					err := h.changePageCallback(ctx, update.CallbackQuery)
					if err != nil {
						slog.Error("telegram error", slog.String("component", "tg_handler"),
							slog.String("op", "changePageCallback"),
							slog.Any("error", err),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
						h.errorCallback(ctx, update.CallbackQuery)
					} else {
						slog.Info("telegram callback processed", slog.String("component", "tg_handler"),
							slog.String("op", "changePageCallback"),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
					}
					continue
				case strings.HasPrefix(update.CallbackQuery.Data, "del_msg"):
					err := h.deleteMessageCallback(ctx, update.CallbackQuery)
					if err != nil {
						slog.Error("telegram error", slog.String("component", "tg_handler"),
							slog.String("op", "deleteMessageCallback"),
							slog.Any("error", err),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
						h.errorCallback(ctx, update.CallbackQuery)
					} else {
						// Отменяем запланированное автоматическое удаление
						h.messageCleaner.CancelDeletion(
							update.CallbackQuery.Message.Chat.ID,
							update.CallbackQuery.Message.MessageID,
						)
						slog.Info("telegram callback processed", slog.String("component", "tg_handler"),
							slog.String("op", "deleteMessageCallback"),
							slog.Group("user",
								slog.Int64("user_id", update.CallbackQuery.From.ID),
								slog.Int64("chat_id", update.CallbackQuery.Message.Chat.ID),
								slog.Int("message_id", update.CallbackQuery.Message.MessageID),
							),
						)
					}
					continue
				}
			}

			if update.Message == nil || update.Message.Text == "" {
				continue
			}

			switch update.Message.Text {
			case "/start":
				h.sendInstructionVideo(ctx, &update)
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
				h.translate(ctx, &update, 1)
				continue
			}
		}
	}
}
