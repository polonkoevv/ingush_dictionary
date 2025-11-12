package tg

import (
	"fmt"
	"test/internal/domain/dictionary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func createDictionaryKeyboard(userDicts []int, dicts []*dictionary.Dictionary) tgbotapi.InlineKeyboardMarkup {
	km := tgbotapi.NewInlineKeyboardMarkup()

	for i := 0; i < len(dicts); i += 2 {
		if i+1 < len(dicts) {
			first_str := "❌ "
			second_str := "❌ "

			first_bd := tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("set_dict_%d", dicts[i].DictID))
			second_bd := tgbotapi.NewInlineKeyboardButtonData(second_str+dicts[i+1].Abbr, fmt.Sprintf("set_dict_%d", dicts[i+1].DictID))
			for _, d := range userDicts {
				if d == dicts[i].DictID {
					first_str = "✅ "
					first_bd = tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("del_dict_%d", dicts[i].DictID))
				}
				if d == dicts[i+1].DictID {
					second_str = "✅ "
					second_bd = tgbotapi.NewInlineKeyboardButtonData(second_str+dicts[i+1].Abbr, fmt.Sprintf("del_dict_%d", dicts[i+1].DictID))

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
			first_bd := tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("set_dict_%d", dicts[i].DictID))
			for _, d := range userDicts {
				if d == dicts[i].DictID {
					first_str = "✅ "
					first_bd = tgbotapi.NewInlineKeyboardButtonData(first_str+dicts[i].Abbr, fmt.Sprintf("del_dict_%d", dicts[i].DictID))
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
