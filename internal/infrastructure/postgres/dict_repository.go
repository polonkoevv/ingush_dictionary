package postgres

import (
	"context"
	"database/sql"
	"test/internal/domain/dictionary"
)

type dictionaryRepository struct {
	db *sql.DB
}

func NewDictionaryRepository(db *sql.DB) dictionary.Repository {
	return &dictionaryRepository{db: db}
}

func (r *dictionaryRepository) GetByID(ctx context.Context, dict_id int) (*dictionary.Dictionary, error) {
	query := `SELECT dict_id, abbr, name, author FROM dictionary WHERE dict_id = $1`

	ds := &dictionary.Dictionary{}

	err := r.db.QueryRowContext(ctx, query, dict_id).Scan(
		&ds.DictID, &ds.Abbr, &ds.Name, &ds.Author,
	)

	if err != nil {
		return nil, err
	}

	return ds, nil
}
func (r *dictionaryRepository) Create(ctx context.Context, dict *dictionary.Dictionary) error {

	query := `INSERT INTO dictionary (abbr, name, author)
	VALUES ($1, $2, $3) RETURNING dict_id
	`

	return r.db.QueryRowContext(ctx, query, dict.Abbr, dict.Name, dict.Author).Scan(&dict.DictID)
}
func (r *dictionaryRepository) Delete(ctx context.Context, dict_id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM dictionary WHERE dict_id = $1`, dict_id)

	return err
}
func (r *dictionaryRepository) Update(ctx context.Context, dict *dictionary.Dictionary) error {

	query := `UPDATE dictionary SET
	abbr = $2, name = $3, author = $4
	WHERE dict_id = $1`

	_, err := r.db.ExecContext(ctx, query, dict.DictID, dict.Abbr, dict.Name, dict.Author)

	return err
}

func (r *dictionaryRepository) GetAll(ctx context.Context) ([]*dictionary.Dictionary, error) {

	query := `SELECT dict_id, abbr, name, author FROM dictionary`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var res []*dictionary.Dictionary

	for rows.Next() {
		var d dictionary.Dictionary

		if err := rows.Scan(&d.DictID, &d.Abbr, &d.Name, &d.Author); err != nil {
			return nil, err
		}

		res = append(res, &d)
	}

	return res, nil
}
