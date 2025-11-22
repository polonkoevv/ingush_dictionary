package tg

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"test/internal/infrastructure/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *TgHandler) setDictCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {

	data := strings.Split(callback.Data, ":")

	dict_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}
	err = h.srv.UserSrv.AddDict(ctx, callback.From.ID, dict_id)
	if err != nil {
		return err
	}
	callbackConfig := tgbotapi.NewCallback(callback.ID, "Словарь добавлен")
	dictKeyboard, err := h.getDictKeyboard(ctx, callback.From.ID)
	editKb := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ChatConfig().ChatID, callback.Message.MessageID, *dictKeyboard)
	h.bot.Send(callbackConfig)
	h.bot.Send(editKb)

	return nil
}

func (h *TgHandler) remDictCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {

	data := strings.Split(callback.Data, ":")

	dict_id, err := strconv.Atoi(data[1])
	if err != nil {
		return err
	}
	err = h.srv.UserSrv.RemoveDict(ctx, callback.From.ID, dict_id)
	if err != nil {
		return err
	}
	callbackConfig := tgbotapi.NewCallback(callback.ID, "Словарь удален")
	dictKeyboard, err := h.getDictKeyboard(ctx, callback.From.ID)
	editKb := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ChatConfig().ChatID, callback.Message.MessageID, *dictKeyboard)
	h.bot.Send(callbackConfig)
	h.bot.Send(editKb)
	return nil
}

func (h *TgHandler) errorCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) {
	callbackConfig := tgbotapi.NewCallback(callback.ID, "Возникла ошибка")
	h.bot.Send(callbackConfig)
}

func (h *TgHandler) deleteMessageCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	err := h.deleteMessageSafe(callback.Message.Chat.ChatConfig().ChatID, callback.Message.MessageID)
	if err != nil {
		return err
	}
	callbackConfig := tgbotapi.NewCallback(callback.ID, "Готово")

	h.bot.Send(callbackConfig)
	return nil
}

func (h *TgHandler) changePageCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {

	// print("MOVING PAGE\n")

	data := strings.Split(callback.Data, ":")
	if len(data) < 4 {
		return fmt.Errorf("invalid pagination callback data: %s", callback.Data)
	}

	wordEscaped := data[1]
	word, err := url.QueryUnescape(wordEscaped)
	if err != nil {
		return err
	}

	language := data[2]
	page_number, err := strconv.Atoi(data[3])
	if err != nil {
		return err
	}
	res, max_quant, err := h.getWord(ctx, word, language, callback.From.ID, page_number)
	if err != nil {
		return err
	}

	if kb := createPaginationKeyboard(word, language, page_number, max_quant, postgres.LIMIT); kb != nil {
		editMsg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			res,
			*kb,
		)
		h.bot.Send(editMsg)
	} else {
		editMsg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, res)
		h.bot.Send(editMsg)
	}
	return nil
}
