package word

import "context"

type Repository interface {
	GetByID(ctx context.Context, word_id int) (*Word, error)
	GetIng(ctx context.Context, word string) ([]*Word, error)
	GetIngFiltered(ctx context.Context, word string, tg_user_id int64) ([]*Word, error)
	GetRus(ctx context.Context, word string) ([]*Word, error)
	GetRusFiltered(ctx context.Context, word string, tg_user_id int64) ([]*Word, error)
	Create(ctx context.Context, word *Word) error
	Delete(ctx context.Context, word_id int) error
	Update(ctx context.Context, word *Word) error
}
