package tg

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *TgHandler) hello(ctx context.Context, update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для перевода слов с русского на ингушский и обратно. Нажмите ‘Сменить язык’ и введите слово.")
	h.srv.UserSrv.CreateOrGetUser(ctx, update)
	h.bot.Send(msg)
}

func (h *TgHandler) help(ctx context.Context, update *tgbotapi.Update) {

	helptext := `📋 *Доступные команды:*


	Перезапустить бота - */start*
	Настройка языка - */language*
	Список словарей - */dictionaries*
	Выбор словарей - */choose*
	Помощь и команды - */help*`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helptext)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}

func (h *TgHandler) deleteMessageSafe(chatID int64, messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := h.bot.Request(deleteMsg)
	return err
}

func (h *TgHandler) changeLanguage(ctx context.Context, update *tgbotapi.Update) {
	language, err := h.srv.UserSrv.ChangeLanguage(ctx, update)
	if err != nil {
		log.Println("change language error:", err)
	}
	switch language {
	case "rus":
		language = "*русский*"
	case "ing":
		language = "*ингушский*"
	}
	text := "Язык оригинала изменен на " + language
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}

func (h *TgHandler) translate(ctx context.Context, update *tgbotapi.Update) {
	language, err := h.srv.UserSrv.GetLanguage(ctx, update)
	if err != nil {
		log.Println("get language error:", err)
		return
	}

	res, _, err := h.getWord(ctx, update.Message.Text, language, update.Message.From.ID)
	if err != nil {
		log.Println("get word error:", err)
		return
	}
	if res == "" {
		nf := "Слово не найдено"
		if language != "rus" { // ing
			nf = "Цу тайпара дош дац"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, nf)
		msg.ReplyToMessageID = update.Message.MessageID
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = update.Message.MessageID
	h.bot.Send(msg)
}

func (h *TgHandler) getWord(ctx context.Context, query, language string, tg_user_id int64) (string, int, error) {

	users_dict, err := h.srv.UserSrv.GetUserDicts(ctx, tg_user_id)

	if err != nil {
		return "", 0, err
	}

	if len(users_dict) == 0 {
		return "Сначала выберите словарь", 0, nil
	}

	words, err := h.srv.WordSrv.GetTranslationFiltered(ctx, query, language, tg_user_id)
	if err != nil {
		return "", 0, err
	}

	quant := len(words)
	res := ""

	if quant == 0 {
		if language == "ing" {
			return "Цу тайпара дош дац", quant, nil
		}
		return "Такое слово не найдено", quant, nil

	}

	switch language {
	case "rus":
		res = "Перевод с русского языка:\n"
		for _, w := range words {
			res += fmt.Sprintf("%s\n%s\n", w.DictAbbr, w.Translation)
			res += "\t" + w.Word + "\n"
		}
	case "ing":
		res = "Перевод с ингушского языка:\n"
		for _, w := range words {
			res += fmt.Sprintf("%s\n%s\n", w.DictAbbr, w.Word)
			res += "\t" + w.Translation + "\n"
		}
	default:
		res = "Непредвиденный язык"
	}

	return res, quant, nil
}

func (h *TgHandler) listDictionaries(ctx context.Context, update *tgbotapi.Update) {
	dicts, err := h.srv.DictSrv.GetAllDictionaries(ctx)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении списка словарей")
		msg.ParseMode = "HTML"
		h.bot.Send(msg)
	}

	msgtext := "*Список словарей:*\n\n\n"

	for i, d := range dicts {

		if d.Author == "" {
			d.Author = "Неизвестен"
		}

		temps := fmt.Sprintf("%d) *%s*\n\n *Автор:* %s\n *Аббревиатура:* %s\n\n\n", i+1, d.Name, d.Author, d.Abbr)
		msgtext += temps
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgtext)
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)

}

func (h *TgHandler) getDictKeyboard(ctx context.Context, tg_user_id int64) (*tgbotapi.InlineKeyboardMarkup, error) {
	userDicts, err := h.srv.UserSrv.GetUserDicts(ctx, tg_user_id)

	if err != nil {
		fmt.Println(err)
		return nil, err
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении спсика словарей")
		// msg.ParseMode = "HTML"
		// h.bot.Send(msg)
	}

	dicts, err := h.srv.DictSrv.GetAllDictionaries(ctx)

	if err != nil {
		return nil, err
		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении спсика словарей")
		// msg.ParseMode = "HTML"
		// h.bot.Send(msg)
	}

	dictKeyboard := createDictionaryKeyboard(userDicts, dicts)
	return &dictKeyboard, nil
}

func (h *TgHandler) chooseDicts(ctx context.Context, update *tgbotapi.Update) {

	dictKeyboard, err := h.getDictKeyboard(ctx, update.Message.From.ID)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении списка словарей")
		fmt.Println(err)
		msg.ParseMode = "HTML"
		h.bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы можете выбрать следующие словари:")
	msg.ReplyMarkup = dictKeyboard
	msg.ParseMode = "Markdown"
	h.bot.Send(msg)
}
