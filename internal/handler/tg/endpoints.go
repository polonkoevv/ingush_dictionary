package tg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"test/internal/infrastructure/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *TgHandler) hello(ctx context.Context, update *tgbotapi.Update) error {

	helloMsg := fmt.Sprintf("Привет, %s! \nЯ бот для перевода слов с русского на ингушский и обратно. \nСначала выбери подходящие тебе словари и выбери язык оригинала.", update.Message.From.UserName)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helloMsg)
	_, err := h.srv.UserSrv.CreateOrGetUser(ctx, update)
	if err != nil {
		return err
	}
	_, err = h.bot.Send(msg)

	return err
}

func (h *TgHandler) help(ctx context.Context, update *tgbotapi.Update) error {

	helptext := `📋 *Доступные команды:*


	Перезапустить бота - */start*
	Настройка языка - */language*
	Список словарей - */dictionaries*
	Выбор словарей - */choose*
	Помощь и команды - */help*`

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, helptext)
	msg.ParseMode = "Markdown"
	_, err := h.bot.Send(msg)

	return err
}

func (h *TgHandler) deleteMessageSafe(chatID int64, messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := h.bot.Request(deleteMsg)
	return err
}

func (h *TgHandler) sendInstructionVideo(ctx context.Context, update *tgbotapi.Update) error {
	video := tgbotapi.NewVideo(update.Message.Chat.ID,
		tgbotapi.FilePath("./assets/instruction.mp4"))
	video.Caption = "📖 Видео-инструкция"
	video.SupportsStreaming = true

	_, err := h.bot.Send(video)
	if err != nil {
		slog.Error("Error while sending instrution")
	}
	return err
}

func (h *TgHandler) changeLanguage(ctx context.Context, update *tgbotapi.Update) error {
	language, err := h.srv.UserSrv.ChangeLanguage(ctx, update)
	if err != nil {
		return err
	}

	var langDisplay string
	switch language {
	case "rus":
		langDisplay = "*русский*"
	case "ing":
		langDisplay = "*ингушский*"
	default:
		return fmt.Errorf("unsupported language: %s", language)
	}

	text := "Язык оригинала изменен на " + langDisplay
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "Markdown"
	_, err = h.bot.Send(msg)
	return err
}

func (h *TgHandler) translate(ctx context.Context, update *tgbotapi.Update, page_number int) error {
	language, err := h.srv.UserSrv.GetLanguage(ctx, update)
	if err != nil {
		return fmt.Errorf("get language error: %w", err)
	}

	word, err := prepareWord(update.Message.Text)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Используются запрещенные символы. Для использования доступны только кириллица, латиница и 1")
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = h.bot.Send(msg)

		return err
	}

	res, max_quant, err := h.getWord(ctx, word, language, update.Message.From.ID, page_number)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = h.bot.Send(msg)

		return err
	}
	if res == "" {
		nf := "Слово не найдено"
		if language != "rus" { // ing
			nf = "Цу тайпара дош дац"
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, nf)
		msg.ReplyToMessageID = update.Message.MessageID
		_, err = h.bot.Send(msg)

		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = update.Message.MessageID
	if km := createPaginationKeyboard(update.Message.Text, language, 1, max_quant, postgres.LIMIT); km != nil {
		msg.ReplyMarkup = km
	}
	_, err = h.bot.Send(msg)

	return err
}

func (h *TgHandler) getWord(ctx context.Context, query, language string, tg_user_id int64, page_number int) (string, int, error) {

	users_dict, err := h.srv.UserSrv.GetUserDicts(ctx, tg_user_id)

	if err != nil {
		return "", 0, err
	}

	if len(users_dict) == 0 {
		return "", 0, errors.New("Сначала выбери словарь")
	}

	query = strings.ToLower(query)

	words, max_quant, err := h.srv.WordSrv.GetTranslationFiltered(ctx, query, language, tg_user_id, page_number)
	if err != nil {
		return "", 0, err
	}

	if len(words) == 0 {
		if language == "ing" {
			return "", max_quant, errors.New("Цу тайпара дош дац")
		}
		return "", max_quant, errors.New("Такое слово не найдено")
	}

	var resBuilder strings.Builder

	switch language {
	case "rus":
		fmt.Fprintf(&resBuilder, "%s РУС -> ИНГ\n\n", strings.ToUpper(query))
		for _, w := range words {
			fmt.Fprintf(&resBuilder, "\n%s\n\t%s\n", w.Translation, w.Word)
		}
	case "ing":
		fmt.Fprintf(&resBuilder, "%s ИНГ –> РУС\n\n", strings.ToUpper(query))
		for _, w := range words {
			fmt.Fprintf(&resBuilder, "\n%s\n\t%s\n", w.Word, w.Translation)
		}
	default:
		return "", max_quant, fmt.Errorf("unsupported language: %s", language)
	}

	return resBuilder.String(), max_quant, nil
}

func (h *TgHandler) listDictionaries(ctx context.Context, update *tgbotapi.Update) error {
	dicts, err := h.srv.DictSrv.GetAllDictionaries(ctx)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении списка словарей")
		msg.ParseMode = "HTML"
		_, err = h.bot.Send(msg)
		return err
	}
	var msgBuilder strings.Builder
	msgBuilder.WriteString("*Список словарей:*\n\n\n")

	for i, d := range dicts {
		author := d.Author
		if author == "" {
			author = "Неизвестен"
		}

		fmt.Fprintf(&msgBuilder, "%d) *%s*\n\n *Автор:* %s\n *Аббревиатура:* %s\n\n\n",
			i+1, d.Name, author, d.Abbr)
	}

	msgtext := msgBuilder.String()

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgtext)
	msg.ParseMode = "Markdown"
	_, err = h.bot.Send(msg)

	return err

}

func (h *TgHandler) getDictKeyboard(ctx context.Context, tg_user_id int64) (*tgbotapi.InlineKeyboardMarkup, error) {
	userDicts, err := h.srv.UserSrv.GetUserDicts(ctx, tg_user_id)

	if err != nil {
		return nil, err
	}

	dicts, err := h.srv.DictSrv.GetAllDictionaries(ctx)

	if err != nil {
		return nil, err
	}

	dictKeyboard := createDictionaryKeyboard(userDicts, dicts)
	return &dictKeyboard, nil
}

func (h *TgHandler) chooseDicts(ctx context.Context, update *tgbotapi.Update) error {

	dictKeyboard, err := h.getDictKeyboard(ctx, update.Message.From.ID)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при получении списка словарей")
		slog.Error("failed to get dictionary keyboard", slog.String("component", "tg_handler"),
			slog.Any("error", err))
		msg.ParseMode = "HTML"
		_, err = h.bot.Send(msg)

		return err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вы можете выбрать следующие словари:")
	msg.ReplyMarkup = dictKeyboard
	msg.ParseMode = "Markdown"
	sentMsg, err := h.bot.Send(msg)
	if err != nil {
		return err
	}

	// Планируем автоматическое удаление через время, указанное в конфиге
	h.messageCleaner.ScheduleDeletion(
		sentMsg.Chat.ID,
		sentMsg.MessageID,
		h.messageTTL,
	)

	return nil
}
