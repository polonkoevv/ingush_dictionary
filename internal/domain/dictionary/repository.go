package dictionary

import "context"

type Repository interface {
	GetByID(ctx context.Context, dict_id int) (*Dictionary, error)
	GetAll(ctx context.Context) ([]*Dictionary, error)
	Create(ctx context.Context, dict *Dictionary) error
	Delete(ctx context.Context, dict_id int) error
	Update(ctx context.Context, dict *Dictionary) error
}
