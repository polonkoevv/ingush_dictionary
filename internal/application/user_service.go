package application

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"test/internal/domain/user"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserService struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateOrGetUser(ctx context.Context, update *tgbotapi.Update) (*user.User, error) {
	// Сначала ищем существующего
	existingUser, err := s.userRepo.GetByTelegramID(ctx, update.Message.From.ID)
	if err != nil && err != sql.ErrNoRows {
		slog.Error("failed to get user by telegram ID",
			slog.String("component", "user_service"),
			slog.String("op", "GetByTelegramID"),
			slog.Int64("tg_user_id", update.Message.From.ID),
			slog.Any("error", err),
		)
		return nil, err
	}

	if existingUser != nil {
		return existingUser, err
	}

	// Создаем нового
	newUser := &user.User{
		TgUserID:     update.Message.From.ID,
		FirstName:    update.Message.From.FirstName,
		Language:     "ing",
		LanguageCode: update.Message.From.LanguageCode,
		SignUpDate:   time.Now(),
	}

	err = s.userRepo.Create(ctx, newUser)
	if err != nil {
		slog.Error("failed to create new user",
			slog.String("component", "user_service"),
			slog.String("op", "Create"),
			slog.Int64("tg_user_id", newUser.TgUserID),
			slog.String("first_name", newUser.FirstName),
			slog.Any("error", err),
		)
		return nil, err
	}

	slog.Info("new user created",
		slog.String("component", "user_service"),
		slog.String("op", "CreateOrGetUser"),
		slog.Int64("tg_user_id", newUser.TgUserID),
		slog.Int64("user_id", newUser.UserID),
		slog.String("first_name", newUser.FirstName),
		slog.String("language_code", newUser.LanguageCode),
		slog.String("default_language", newUser.Language),
	)

	return newUser, err
}

func (s *UserService) ChangeLanguage(ctx context.Context, update *tgbotapi.Update) (string, error) {
	user, err := s.CreateOrGetUser(ctx, update)
	if err != nil {
		slog.Error("failed to get user for language change",
			slog.String("component", "user_service"),
			slog.String("op", "ChangeLanguage"),
			slog.Int64("tg_user_id", update.Message.From.ID),
			slog.Any("error", err),
		)
		return "", err
	}
	if user == nil {
		err := fmt.Errorf("user not found")
		slog.Error("user not found for language change",
			slog.String("component", "user_service"),
			slog.String("op", "ChangeLanguage"),
			slog.Int64("tg_user_id", update.Message.From.ID),
		)
		return "", err
	}

	oldLang := user.Language
	newLang := "rus"
	if user.Language == "rus" {
		newLang = "ing"
	}

	err = s.userRepo.UpdateLanguage(ctx, update.Message.From.ID, newLang)
	if err != nil {
		slog.Error("failed to update user language",
			slog.String("component", "user_service"),
			slog.String("op", "ChangeLanguage"),
			slog.Int64("tg_user_id", update.Message.From.ID),
			slog.String("old_language", oldLang),
			slog.String("new_language", newLang),
			slog.Any("error", err),
		)
		return "", err
	}

	slog.Info("user language changed",
		slog.String("component", "user_service"),
		slog.String("op", "ChangeLanguage"),
		slog.Int64("tg_user_id", update.Message.From.ID),
		slog.String("old_language", oldLang),
		slog.String("new_language", newLang),
	)

	return newLang, err
}

func (s *UserService) GetLanguage(ctx context.Context, update *tgbotapi.Update) (string, error) {

	user, err := s.CreateOrGetUser(ctx, update)
	if err != nil {
		slog.Error("failed to get user language",
			slog.String("component", "user_service"),
			slog.String("op", "GetLanguage"),
			slog.Int64("tg_user_id", update.Message.From.ID),
			slog.Any("error", err),
		)
		return "", err
	}
	if user == nil {
		err := fmt.Errorf("user not found")
		slog.Error("user not found when getting language",
			slog.String("component", "user_service"),
			slog.String("op", "GetLanguage"),
			slog.Int64("tg_user_id", update.Message.From.ID),
		)
		return "", err
	}
	return user.Language, nil
}

func (s *UserService) GetUserDicts(ctx context.Context, tg_user_id int64) ([]int, error) {
	res, err := s.userRepo.GetUserDicts(ctx, tg_user_id)

	if err != nil {
		slog.Error("failed to get user dictionaries",
			slog.String("component", "user_service"),
			slog.String("op", "GetUserDicts"),
			slog.Int64("tg_user_id", tg_user_id),
			slog.Any("error", err),
		)
		return nil, err
	}

	slog.Debug("user dictionaries retrieved",
		slog.String("component", "user_service"),
		slog.String("op", "GetUserDicts"),
		slog.Int64("tg_user_id", tg_user_id),
		slog.Int("dicts_count", len(res)),
	)

	return res, nil
}

func (s *UserService) AddDict(ctx context.Context, tg_user_id int64, dict_id int) error {
	err := s.userRepo.AddDict(ctx, tg_user_id, dict_id)

	if err != nil {
		slog.Error("failed to add dictionary to user",
			slog.String("component", "user_service"),
			slog.String("op", "AddDict"),
			slog.Int64("tg_user_id", tg_user_id),
			slog.Int("dict_id", dict_id),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("dictionary added to user",
		slog.String("component", "user_service"),
		slog.String("op", "AddDict"),
		slog.Int64("tg_user_id", tg_user_id),
		slog.Int("dict_id", dict_id),
	)

	return err
}

func (s *UserService) RemoveDict(ctx context.Context, tg_user_id int64, dict_id int) error {
	err := s.userRepo.RemoveDict(ctx, tg_user_id, dict_id)

	if err != nil {
		slog.Error("failed to remove dictionary from user",
			slog.String("component", "user_service"),
			slog.String("op", "RemoveDict"),
			slog.Int64("tg_user_id", tg_user_id),
			slog.Int("dict_id", dict_id),
			slog.Any("error", err),
		)
		return err
	}

	slog.Info("dictionary removed from user",
		slog.String("component", "user_service"),
		slog.String("op", "RemoveDict"),
		slog.Int64("tg_user_id", tg_user_id),
		slog.Int("dict_id", dict_id),
	)

	return err
}
