package user

import "time"

type User struct {
	UserID       int64     `json:"user_id" db:"user_id"`
	TgUserID     int64     `json:"tg_user_id" db:"tg_user_id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LanguageCode string    `json:"language_code" db:"language_code"`
	SignUpDate   time.Time `json:"signup_date" db:"signup_date"`
	Language     string    `json:"language" db:"language"`
}
