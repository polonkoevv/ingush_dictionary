package tg

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MessageToDelete хранит информацию о сообщении, которое нужно удалить
type MessageToDelete struct {
	ChatID    int64
	MessageID int
	DeleteAt  time.Time
}

// MessageCleaner управляет автоматическим удалением сообщений
type MessageCleaner struct {
	bot      *tgbotapi.BotAPI
	messages map[string]*MessageToDelete // ключ: "chatID:messageID"
	mu       sync.RWMutex
	stopChan chan struct{}
}

// NewMessageCleaner создает новый экземпляр MessageCleaner
func NewMessageCleaner(bot *tgbotapi.BotAPI) *MessageCleaner {
	return &MessageCleaner{
		bot:      bot,
		messages: make(map[string]*MessageToDelete),
		stopChan: make(chan struct{}),
	}
}

// ScheduleDeletion планирует удаление сообщения через указанное время
// ttl - время жизни сообщения (например, 24 часа)
func (mc *MessageCleaner) ScheduleDeletion(chatID int64, messageID int, ttl time.Duration) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.messageKey(chatID, messageID)
	mc.messages[key] = &MessageToDelete{
		ChatID:    chatID,
		MessageID: messageID,
		DeleteAt:  time.Now().Add(ttl),
	}

	slog.Debug("scheduled message deletion",
		slog.String("component", "tg_handler"),
		slog.String("op", "ScheduleDeletion"),
		slog.Int64("chat_id", chatID),
		slog.Int("message_id", messageID),
		slog.String("delete_at", time.Now().Add(ttl).Format(time.RFC3339)),
	)
}

// CancelDeletion отменяет запланированное удаление сообщения
func (mc *MessageCleaner) CancelDeletion(chatID int64, messageID int) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	key := mc.messageKey(chatID, messageID)
	if _, exists := mc.messages[key]; exists {
		delete(mc.messages, key)
		slog.Debug("cancelled message deletion",
			slog.String("component", "tg_handler"),
			slog.String("op", "CancelDeletion"),
			slog.Int64("chat_id", chatID),
			slog.Int("message_id", messageID),
		)
	}
}

// Start запускает фоновую горутину для проверки и удаления сообщений
func (mc *MessageCleaner) Start(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute) // проверяем каждую минуту
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("message cleaner stopped",
				slog.String("component", "tg_handler"),
				slog.String("op", "MessageCleaner.Start"))
			return
		case <-mc.stopChan:
			return
		case <-ticker.C:
			mc.processDeletions()
		}
	}
}

// processDeletions проверяет и удаляет сообщения, время которых истекло
func (mc *MessageCleaner) processDeletions() {
	mc.mu.Lock()
	now := time.Now()
	var toDelete []*MessageToDelete

	for key, msg := range mc.messages {
		if now.After(msg.DeleteAt) {
			toDelete = append(toDelete, msg)
			delete(mc.messages, key)
		}
	}
	mc.mu.Unlock()

	// Удаляем сообщения вне блокировки
	for _, msg := range toDelete {
		mc.deleteMessage(msg.ChatID, msg.MessageID)
	}
}

// deleteMessage безопасно удаляет сообщение
func (mc *MessageCleaner) deleteMessage(chatID int64, messageID int) {
	deleteMsg := tgbotapi.NewDeleteMessage(chatID, messageID)
	_, err := mc.bot.Request(deleteMsg)
	if err != nil {
		// Игнорируем ошибку, если сообщение уже удалено или нельзя удалить
		slog.Debug("failed to auto-delete message (may be already deleted)",
			slog.String("component", "tg_handler"),
			slog.String("op", "deleteMessage"),
			slog.Int64("chat_id", chatID),
			slog.Int("message_id", messageID),
			slog.Any("error", err),
		)
	} else {
		slog.Info("auto-deleted message",
			slog.String("component", "tg_handler"),
			slog.String("op", "deleteMessage"),
			slog.Int64("chat_id", chatID),
			slog.Int("message_id", messageID),
		)
	}
}

// messageKey создает ключ для map из chatID и messageID
func (mc *MessageCleaner) messageKey(chatID int64, messageID int) string {
	return fmt.Sprintf("%d:%d", chatID, messageID)
}

// Stop останавливает cleaner
func (mc *MessageCleaner) Stop() {
	close(mc.stopChan)
}

