package tg

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *TgHandler) addDictCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	dict_id, err := strconv.Atoi(callback.Data[9:])
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = h.srv.UserSrv.AddDict(ctx, callback.From.ID, dict_id)
	if err != nil {
		fmt.Println(err)
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
	dict_id, err := strconv.Atoi(callback.Data[9:])
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = h.srv.UserSrv.RemoveDict(ctx, callback.From.ID, dict_id)
	if err != nil {
		fmt.Println(err)
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
