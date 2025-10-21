package tg

import (
	"fmt"
	"log"
	"strings"
	"test/internal/handler/tg/utils"
	"test/internal/service"
	"test/internal/storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgHandler struct {
	bot *tgbotapi.BotAPI
	s   storage.Storage
	svc service.LanguageService
}

func NewTgHandler(bot *tgbotapi.BotAPI, s storage.Storage, svc service.LanguageService) *TgHandler {
	return &TgHandler{bot: bot, s: s, svc: svc}
}

func (h *TgHandler) ChangeLanguage(chatID int64) (string, error) {
	return h.s.ChangeLanguage(chatID)
}

func (h *TgHandler) GetLanguage(chatID int64) (string, error) {
	return h.s.GetLanguage(chatID)
}

// Run запускает цикл обработки апдейтов TG, включая сообщения и callback-кнопки.
func (h *TgHandler) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := h.bot.GetUpdatesChan(u)

	replyKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сменить язык"),
		),
	)

	for update := range updates {
		if update.CallbackQuery != nil {
			h.handleCallback(update)
			continue
		}

		if update.InlineQuery != nil {
			query := update.InlineQuery.Query
			var results []interface{}

			if query == "" {
				// Если запрос пустой - показываем примеры использования
				results = []interface{}{
					tgbotapi.NewInlineQueryResultArticleHTML(
						"1",
						"Инструкция",
						"Введите слово для перевода, Вам будет предложен список переводов",
					),
				}
			} else {

				log.Println("QUERY:", query)

				ing, _, err := h.getWord(query, "ing")
				rus, _, err := h.getWord(query, "rus")
				if ing == "" {
					ing = "Цу тайпара дош дац"
				}
				if rus == "" {
					rus = "Слово не найдено"
				}

				log.Println("ING:", ing)
				log.Println("RUS:", rus)

				// Пытаемся распарсить число из запроса
				if err == nil {
					results = []interface{}{
						tgbotapi.NewInlineQueryResultArticleHTML(
							"1",
							"Ингушский",
							ing,
						),
						tgbotapi.NewInlineQueryResultArticleHTML(
							"2",
							"Русский",
							rus,
						),
					}
				} else {
					results = []interface{}{
						tgbotapi.NewInlineQueryResultArticleHTML(
							"1",
							"Ошибка",
							"Пожалуйста, введите число (например: <code>50</code>)",
						),
					}
				}
			}

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				Results:       results,
				CacheTime:     1, // секунды
			}

			if _, err := h.bot.Request(inlineConfig); err != nil {
				log.Printf("Error sending inline response: %v", err)
			}
		}

		if update.Message == nil {
			continue
		}

		if update.Message.Text == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный тип сообщения")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = replyKeyboard
			h.bot.Send(msg)
			continue
		}

		if update.Message.Text == "/start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для перевода слов с русского на ингушский и обратно. Нажмите ‘Сменить язык’ и введите слово.")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = replyKeyboard
			h.bot.Send(msg)
			continue
		}

		if update.Message.Text == "Сменить язык" || update.Message.Text == "Русский | Сменить язык" || update.Message.Text == "Ингушский | Сменить язык" {
			language, err := h.ChangeLanguage(update.Message.Chat.ID)
			if err != nil {
				log.Println("change language error:", err)
			}
			text := "Выбранный язык: русский"
			replyKeyboard = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Русский | Сменить язык"),
				),
			)
			if language == "ing" {
				text = "Выбранный язык: ингушский"
				replyKeyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Ингушский | Сменить язык"),
					),
				)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = replyKeyboard
			h.bot.Send(msg)
			continue
		}

		// Обычный текст — перевод
		language, err := h.svc.GetLanguage(update.Message.Chat.ID)
		if err != nil {
			log.Println("get language error:", err)
			continue
		}

		res, total, err := h.getWord(update.Message.Text, language)
		if err != nil {
			log.Println("get word error:", err)
			continue
		}
		if res == "" {
			nf := "Слово не найдено"
			if language != "rus" { // ing
				nf = "Цу тайпара дош дац"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, nf)
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = replyKeyboard
			h.bot.Send(msg)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, res)
		msg.ParseMode = "HTML"
		msg.ReplyToMessageID = update.Message.MessageID
		kb := utils.GetButtons(total, 0, 7, language, update.Message.Text)
		msg.ReplyMarkup = kb
		h.bot.Send(msg)
	}
}

func (h *TgHandler) getWord(text, language string) (string, int, error) {
	if language == "rus" {
		return h.svc.RusToIng(text)
	}
	return h.svc.IngToRus(text)
}

func (h *TgHandler) handleCallback(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	// nav:prev|next:page:per:total:from:word
	// sel:globalIndex:page:per:total:from:word
	parts := strings.Split(data, ":")
	if len(parts) < 2 {
		return
	}

	chatID := update.CallbackQuery.Message.Chat.ID
	msgID := update.CallbackQuery.Message.MessageID

	switch parts[0] {
	case "nav":
		if len(parts) < 7 {
			return
		}
		dir := parts[1]
		page := atoiSafe(parts[2])
		per := atoiSafe(parts[3])
		total := atoiSafe(parts[4])
		from := parts[5]
		word := parts[6]
		if dir == "next" {
			page++
		} else {
			if page > 0 {
				page--
			}
		}
		kb := utils.GetButtons(total, page, per, from, word)
		edit := tgbotapi.NewEditMessageReplyMarkup(chatID, msgID, *kb)
		h.bot.Send(edit)

	case "sel":
		// Сейчас просто подтверждаем выбор значения. Можно расширить деталями.
		if len(parts) < 7 {
			return
		}
		idx := atoiSafe(parts[1])
		word := parts[6]
		text := fmt.Sprintf("<b>%s</b> — выбран вариант #%d", escapeHTML(word), idx)
		edit := tgbotapi.NewEditMessageText(chatID, msgID, text)
		edit.ParseMode = "HTML"
		h.bot.Send(edit)
	}

	// обязательно ответим callback, чтобы убрать «часики»
	h.bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
}

func atoiSafe(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return n
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func escapeHTML(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", "\"", "&quot;")
	return r.Replace(s)
}
