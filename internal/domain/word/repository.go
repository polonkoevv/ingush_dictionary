package word

import "context"

type Repository interface {
	GetByID(ctx context.Context, word_id int) (*Word, error)
	GetIng(ctx context.Context, word string, page_number int) ([]*Word, int, error)
	GetIngFiltered(ctx context.Context, word string, tg_user_id int64, page_number int) ([]*Word, int, error)
	GetRus(ctx context.Context, word string, page_number int) ([]*Word, int, error)
	GetRusFiltered(ctx context.Context, word string, tg_user_id int64, page_number int) ([]*Word, int, error)
	Create(ctx context.Context, word *Word) error
	Delete(ctx context.Context, word_id int) error
	Update(ctx context.Context, word *Word) error
}
