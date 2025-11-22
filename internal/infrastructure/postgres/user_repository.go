package postgres

import (
	"context"
	"database/sql"
	"test/internal/domain/user"

	_ "github.com/lib/pq"
)

// type Repository interface {
// 	GetByTelegramID(ctx context.Context, tg_user_id int) (User, error)
// 	Create(ctx context.Context, user User) error
// 	Delete(ctx context.Context, tg_user_id int) error
// 	Update(ctx context.Context, user User) error
// 	UpdateLanguage(ctx context.Context, telegramID int64, language string) error
// }

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*user.User, error) {

	query := `SELECT user_id, tg_user_id, first_name, language_code, signup_date, language
            FROM users WHERE tg_user_id = $1
    `

	u := &user.User{}

	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(&u.UserID, &u.TgUserID, &u.FirstName,
		&u.LanguageCode, &u.SignUpDate, &u.Language)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *userRepository) Create(ctx context.Context, user *user.User) error {

	query := `INSERT INTO users (tg_user_id, first_name, language_code, signup_date, language) VALUES ($1, $2, $3, $4, $5) RETURNING user_id`

	err := r.db.QueryRowContext(ctx, query, user.TgUserID, user.FirstName, user.LanguageCode, user.SignUpDate, user.Language).Scan(&user.UserID)

	return err
}
func (r *userRepository) Delete(ctx context.Context, tg_user_id int64) error {
	query := `DELETE FROM users WHERE tg_user_id = $1`

	_, err := r.db.ExecContext(ctx, query, tg_user_id)

	return err
}
func (r *userRepository) Update(ctx context.Context, user *user.User) error {
	query := `UPDATE users SET first_name = $2, language_code = $3, language = $4 WHERE tg_user_id = $1`

	_, err := r.db.ExecContext(ctx, query, user.TgUserID, user.FirstName, user.LanguageCode, user.Language)

	return err
}
func (r *userRepository) UpdateLanguage(ctx context.Context, telegramID int64, language string) error {

	query := `UPDATE users SET language = $2 WHERE tg_user_id = $1`

	_, err := r.db.ExecContext(ctx, query, telegramID, language)

	return err
}

func (r *userRepository) AddDict(ctx context.Context, tg_user_id int64, dict_id int) error {

	query := `INSERT INTO users_to_dict (tg_user_id, dict_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, tg_user_id, dict_id)

	return err
}

func (r *userRepository) RemoveDict(ctx context.Context, tg_user_id int64, dict_id int) error {
	query := `DELETE FROM users_to_dict WHERE tg_user_id = $1 AND dict_id = $2`

	_, err := r.db.ExecContext(ctx, query, tg_user_id, dict_id)

	return err
}

func (r *userRepository) GetUserDicts(ctx context.Context, tg_user_id int64) ([]int, error) {
	query := `SELECT dict_id FROM users_to_dict WHERE tg_user_id = $1`

	var res []int

	rows, err := r.db.QueryContext(ctx, query, tg_user_id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t int

		if err := rows.Scan(&t); err != nil {
			return nil, err
		}

		res = append(res, t)
	}

	return res, rows.Err()
}
