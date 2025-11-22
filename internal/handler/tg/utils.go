package tg

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"test/internal/domain/dictionary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createDictionaryKeyboard(userDicts []int, dicts []*dictionary.Dictionary) tgbotapi.InlineKeyboardMarkup {
	km := tgbotapi.NewInlineKeyboardMarkup()

	for i := 0; i < len(dicts); i += 2 {
		if i+1 < len(dicts) {
			first_str := "❌ "
			second_str := "❌ "

			first_bd := tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("set_dict:%d", dicts[i].DictID))
			second_bd := tgbotapi.NewInlineKeyboardButtonData(second_str+dicts[i+1].Abbr, fmt.Sprintf("set_dict:%d", dicts[i+1].DictID))
			for _, d := range userDicts {
				if d == dicts[i].DictID {
					first_str = "✅ "
					first_bd = tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("rem_dict:%d", dicts[i].DictID))
				}
				if d == dicts[i+1].DictID {
					second_str = "✅ "
					second_bd = tgbotapi.NewInlineKeyboardButtonData(second_str+dicts[i+1].Abbr, fmt.Sprintf("rem_dict:%d", dicts[i+1].DictID))

				}

				if first_str == "✅ " && second_str == "✅ " {
					break
				}
			}

			km.InlineKeyboard = append(km.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					first_bd,
					second_bd,
				),
			)
		} else {
			first_str := "❌ "
			first_bd := tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("set_dict:%d", dicts[i].DictID))
			for _, d := range userDicts {
				if d == dicts[i].DictID {
					first_str = "✅ "
					first_bd = tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("rem_dict:%d", dicts[i].DictID))
					break
				}
			}
			km.InlineKeyboard = append(km.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					first_bd,
				),
			)
		}

	}

	km.InlineKeyboard = append(km.InlineKeyboard,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Готово", "del_msg"),
		),
	)
	return km
}

const paginationButtonsCount = 5

func createPaginationKeyboard(word string, language string, pageNumber int, quantity int, limit int) *tgbotapi.InlineKeyboardMarkup {
	km := tgbotapi.NewInlineKeyboardMarkup()

	if limit <= 0 || quantity <= 0 {
		return nil
	}

	totalPages := int(math.Ceil(float64(quantity) / float64(limit)))
	if totalPages <= 1 {
		return nil
	}

	if pageNumber < 1 {
		pageNumber = 1
	} else if pageNumber > totalPages {
		pageNumber = totalPages
	}

	windowStart := pageNumber - paginationButtonsCount/2
	if windowStart < 1 {
		windowStart = 1
	}
	windowEnd := windowStart + paginationButtonsCount - 1
	if windowEnd > totalPages {
		windowEnd = totalPages
		windowStart = windowEnd - paginationButtonsCount + 1
		if windowStart < 1 {
			windowStart = 1
		}
	}

	wordEscaped := url.QueryEscape(word)

	numRow := tgbotapi.NewInlineKeyboardRow()
	for page := windowStart; page <= windowEnd; page++ {
		label := fmt.Sprintf("%d", page)
		if page == pageNumber {
			label = fmt.Sprintf("[%d]", page)
		}
		numRow = append(numRow,
			tgbotapi.NewInlineKeyboardButtonData(
				label,
				fmt.Sprintf("change_page:%s:%s:%d", wordEscaped, language, page),
			),
		)
	}
	km.InlineKeyboard = append(km.InlineKeyboard, numRow)

	moveRow := tgbotapi.NewInlineKeyboardRow()
	if pageNumber > 1 {
		moveRow = append(moveRow,
			tgbotapi.NewInlineKeyboardButtonData(
				"◀️ Назад",
				fmt.Sprintf("change_page:%s:%s:%d", wordEscaped, language, pageNumber-1),
			),
		)
	}
	if pageNumber < totalPages {
		moveRow = append(moveRow,
			tgbotapi.NewInlineKeyboardButtonData(
				"Вперёд ▶️",
				fmt.Sprintf("change_page:%s:%s:%d", wordEscaped, language, pageNumber+1),
			),
		)
	}
	if len(moveRow) > 0 {
		km.InlineKeyboard = append(km.InlineKeyboard, moveRow)
	}

	return &km
}

const allowedChars = " 1АаБбВвГгДдЕеЖжЗзИиЙйКкЛлМмНнОоПпРрСсТтУуФфХхЦцЧчШшЩщЪъЫыЬьЭэЮюЯяAaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"

func containsOnlyChars(input, allowedChars string) bool {
	// Создаем map для быстрой проверки
	allowed := make(map[rune]bool)
	for _, char := range allowedChars {
		allowed[char] = true
	}

	// Проверяем каждый символ входной строки
	for _, char := range input {
		if !allowed[char] {
			return false
		}
	}
	return true
}

func prepareWord(word string) (string, error) {

	if !containsOnlyChars(word, allowedChars) {
		return "", errors.New("bad characters")
	}

	word = strings.ToLower(
		strings.Trim(word, ""),
	)

	return word, nil
}
