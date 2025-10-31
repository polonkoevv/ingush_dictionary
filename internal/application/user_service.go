package application

import (
	"context"
	"database/sql"
	"fmt"
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
	return newUser, err
}

func (s *UserService) ChangeLanguage(ctx context.Context, update *tgbotapi.Update) (string, error) {
	user, err := s.CreateOrGetUser(ctx, update)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	newLang := "rus"
	if user.Language == "rus" {
		newLang = "ing"
	}

	err = s.userRepo.UpdateLanguage(ctx, update.Message.From.ID, newLang)
	return newLang, err
}

func (s *UserService) GetLanguage(ctx context.Context, update *tgbotapi.Update) (string, error) {

	user, err := s.CreateOrGetUser(ctx, update)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}
	return user.Language, nil
}

func (s *UserService) GetUserDicts(ctx context.Context, tg_user_id int64) ([]int, error) {
	res, err := s.userRepo.GetUserDicts(ctx, tg_user_id)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserService) AddDict(ctx context.Context, tg_user_id int64, dict_id int) error {
	err := s.userRepo.AddDict(ctx, tg_user_id, dict_id)

	return err
}

func (s *UserService) RemoveDict(ctx context.Context, tg_user_id int64, dict_id int) error {
	err := s.userRepo.RemoveDict(ctx, tg_user_id, dict_id)

	return err
}
