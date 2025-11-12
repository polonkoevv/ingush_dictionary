package postgres

import (
	"context"
	"database/sql"
	"strings"
	"test/internal/domain/word"
)

// type Repository interface {
// 	GetByID(ctx context.Context, word_id int) (*Word, error)
// 	GetIng(ctx context.Context, word string) ([]*Word, error)
// 	GetRus(ctx context.Context, word string) ([]*Word, error)
// 	Create(ctx context.Context, word *Word) error
// 	Delete(ctx context.Context, word_id int) error
// 	Update(ctx context.Context, word *Word) error
// }

type wordRepository struct {
	db *sql.DB
}

func NewWordRepository(db *sql.DB) word.Repository {
	return &wordRepository{db: db}
}

func (r *wordRepository) GetByID(ctx context.Context, word_id int) (*word.Word, error) {

	query := `SELECT word_id, word, translation, speech_part, topic, dict_id FROM word WHERE word_id = $1 `

	w := &word.Word{}

	err := r.db.QueryRowContext(ctx, query, word_id).Scan(&w.WordID, &w.Word,
		&w.Translation, &w.SpeechPart, &w.Topic, &w.DictID)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return w, err
}

func (r *wordRepository) GetIng(ctx context.Context, query string) ([]*word.Word, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*word.Word{}, nil
	}

	const sql = `
        SELECT word_id, word, translation, speech_part, topic, dict_id
        FROM word
        WHERE word ILIKE $1
        ORDER BY
            CASE
                WHEN word = $2 THEN 1
                WHEN word ILIKE $3 THEN 2
                ELSE 3
            END,
            translation ASC
        LIMIT 20
    `

	rows, err := r.db.QueryContext(ctx, sql, "%"+query+"%", query, query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*word.Word
	for rows.Next() {
		var w word.Word
		if err := rows.Scan(&w.WordID, &w.Word, &w.Translation, &w.SpeechPart, &w.Topic, &w.DictID); err != nil {
			return nil, err
		}
		results = append(results, &w)
	}
	return results, rows.Err()
}

func (r *wordRepository) GetIngFiltered(ctx context.Context, query string, tg_user_id int64) ([]*word.Word, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*word.Word{}, nil
	}

	const sql = `
	SELECT w.word_id, w.word, w.translation, w.speech_part, w.topic, w.dict_id, d.abbr
        FROM word w
        JOIN dictionary d ON d.dict_id = w.dict_id
        WHERE word ILIKE $1 AND w.dict_id = ANY(SELECT dict_id FROM users_to_dict WHERE tg_user_id = $4)
        ORDER BY
            CASE
                WHEN word = $2 THEN 1
                WHEN word ILIKE $3 THEN 2
                ELSE 3
            END,
            translation ASC
        LIMIT 20
    `

	rows, err := r.db.QueryContext(ctx, sql, "%"+query+"%", query, query+"%", tg_user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*word.Word
	for rows.Next() {
		var w word.Word
		if err := rows.Scan(&w.WordID, &w.Word, &w.Translation, &w.SpeechPart, &w.Topic, &w.DictID, &w.DictAbbr); err != nil {
			return nil, err
		}
		results = append(results, &w)
	}
	return results, rows.Err()
}

func (r *wordRepository) GetRus(ctx context.Context, query string) ([]*word.Word, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*word.Word{}, nil
	}

	const sql = `
        SELECT word_id, word, translation, speech_part, topic, dict_id
        FROM word
        WHERE translation ILIKE $1
        ORDER BY
            CASE
                WHEN translation = $2 THEN 1
                WHEN translation ILIKE $3 THEN 2
                ELSE 3
            END,
            word ASC
        LIMIT 20
    `

	rows, err := r.db.QueryContext(ctx, sql, "%"+query+"%", query, query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*word.Word
	for rows.Next() {
		var w word.Word
		if err := rows.Scan(&w.WordID, &w.Word, &w.Translation, &w.SpeechPart, &w.Topic, &w.DictID); err != nil {
			return nil, err
		}
		results = append(results, &w)
	}
	return results, rows.Err()
}

func (r *wordRepository) GetRusFiltered(ctx context.Context, query string, tg_user_id int64) ([]*word.Word, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []*word.Word{}, nil
	}

	const sql = `
        SELECT w.word_id, w.word, w.translation, w.speech_part, w.topic, w.dict_id, d.abbr
        FROM word w
        JOIN dictionary d ON d.dict_id = w.dict_id
        WHERE translation ILIKE $1 AND w.dict_id = ANY(SELECT dict_id FROM users_to_dict WHERE tg_user_id = $4)
        ORDER BY
            CASE
                WHEN translation = $2 THEN 1
                WHEN translation ILIKE $3 THEN 2
                ELSE 3
            END,
            word ASC
        LIMIT 20
    `

	rows, err := r.db.QueryContext(ctx, sql, "%"+query+"%", query, query+"%", tg_user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*word.Word
	for rows.Next() {
		var w word.Word
		if err := rows.Scan(&w.WordID, &w.Word, &w.Translation, &w.SpeechPart, &w.Topic, &w.DictID, &w.DictAbbr); err != nil {
			return nil, err
		}
		results = append(results, &w)
	}
	return results, rows.Err()
}

func (r *wordRepository) Create(ctx context.Context, word *word.Word) error {
	query := `INSERT INTO word (word, translation, speech_part, topic, dict_id)
	VALUES ($1, $2, $3, $4, $5) RETURNING word_id`

	err := r.db.QueryRowContext(ctx, query, word.Word, word.Translation,
		word.SpeechPart, word.Topic, word.DictID).Scan(&word.WordID)

	return err
}
func (r *wordRepository) Delete(ctx context.Context, word_id int) error {
	query := `DELETE FROM word WHERE word_id = $1`

	_, err := r.db.ExecContext(ctx, query, word_id)

	return err
}
func (r *wordRepository) Update(ctx context.Context, word *word.Word) error {

	query := `UPDATE word SET word = $2, translation = $3, speech_part = $4, topic = $5, dict_id = $6 WHERE word_id = $1`

	_, err := r.db.ExecContext(ctx, query, word.WordID, word.Word, word.Translation, word.SpeechPart, word.Topic, word.DictID)

	return err
}
