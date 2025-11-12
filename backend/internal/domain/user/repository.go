package user

import "context"

type Repository interface {
	GetByTelegramID(ctx context.Context, tg_user_id int64) (*User, error)
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, tg_user_id int64) error
	Update(ctx context.Context, user *User) error
	UpdateLanguage(ctx context.Context, telegramID int64, language string) error

	AddDict(ctx context.Context, tg_user_id int64, dict_id int) error
	RemoveDict(ctx context.Context, tg_user_id int64, dict_id int) error
	GetUserDicts(ctx context.Context, tg_user_id int64) ([]int, error)
}
