package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetButtons(num int, total int, from, word string) *tgbotapi.InlineKeyboardMarkup {

	// total num page index from word

	buttons := make([][]tgbotapi.InlineKeyboardButton, 2)
	buttons[1] = []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("<", "trans_prev"),
		tgbotapi.NewInlineKeyboardButtonData(">", "trans_next"),
	}

	buttons[0] = make([]tgbotapi.InlineKeyboardButton, num)
	for i := 0; i < num; i++ {
		buttons[0][i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+1), fmt.Sprintf("%d:%d:%d:%d:%s:%s", total, num, i+1, i, from, word))
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons...)

	return &keyboard
}
