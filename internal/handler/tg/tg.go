package tg

import (
	"test/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgHandler struct {
	bot *tgbotapi.BotAPI
	s   storage.Storage
}

func NewTgHandler(bot *tgbotapi.BotAPI, s storage.Storage) *TgHandler {
	return &TgHandler{bot: bot, s: s}
}

func (h *TgHandler) ChangeLanguage(chatID int64) (string, error) {
	return h.s.ChangeLanguage(chatID)
}

func (h *TgHandler) GetLanguage(chatID int64) (string, error) {
	return h.s.GetLanguage(chatID)
}
