package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetButtons создает инлайн-клавиатуру с пагинацией.
// total — общее число значений, page — текущая страница (0-based), per — элементов на странице,
// from — направление перевода ("rus"|"ing"), word — исходный запрос.
func GetButtons(total int, page int, per int, from, word string) *tgbotapi.InlineKeyboardMarkup {
	if per <= 0 {
		per = 7
	}

	// Рассчитаем сколько элементов показывать на текущей странице
	start := page * per
	if start < 0 {
		start = 0
	}
	remaining := total - start
	if remaining < 0 {
		remaining = 0
	}
	count := per
	if remaining < per {
		count = remaining
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 2)

	// Ряд с номерами значений на странице (1..count)
	rows[0] = make([]tgbotapi.InlineKeyboardButton, count)
	for i := 0; i < count; i++ {
		globalIndex := start + i + 1
		data := fmt.Sprintf("sel:%d:%d:%d:%d:%s:%s", globalIndex, page, per, total, from, word)
		rows[0][i] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i+1), data)
	}

	// Ряд навигации
	prevData := fmt.Sprintf("nav:prev:%d:%d:%d:%s:%s", page, per, total, from, word)
	nextData := fmt.Sprintf("nav:next:%d:%d:%d:%s:%s", page, per, total, from, word)
	rows[1] = []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("<", prevData),
		tgbotapi.NewInlineKeyboardButtonData(">", nextData),
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard
}
