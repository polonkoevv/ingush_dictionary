package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"test/internal/config"
	"test/internal/domain/dictionary"
	"test/internal/domain/user"
	"test/internal/domain/word"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	UserRep user.Repository
	WordRep word.Repository
	DictRep dictionary.Repository
}

func CreateRepositories(db *sql.DB) (Repository, error) {
	usrRep := NewUserRepository(db)
	wrdRep := NewWordRepository(db)
	dictRep := NewDictionaryRepository(db)

	return Repository{
		UserRep: usrRep,
		WordRep: wrdRep,
		DictRep: dictRep,
	}, nil
}

func NewPostgresConnection(cfg *config.Config) (*sql.DB, error) {
	// Получаем строку подключения из переменных окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := cfg.Postgres.Host
		port := cfg.Postgres.Port
		user := cfg.Postgres.User
		password := cfg.Postgres.Password
		dbname := cfg.Postgres.DBName

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
