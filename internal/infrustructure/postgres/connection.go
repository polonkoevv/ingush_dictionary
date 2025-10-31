package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"test/internal/domain/dictionary"
	"test/internal/domain/user"
	"test/internal/domain/word"
	"time"
)

func CreateRepositories(db *sql.DB) (user.Repository, word.Repository, dictionary.Repository, error) {
	usrRep := NewUserRepository(db)
	wrdRep := NewWordRepository(db)
	dictRep := NewDictionaryRepository(db)

	return usrRep, wrdRep, dictRep, nil
}

func NewPostgresConnection() (*sql.DB, error) {
	// Получаем строку подключения из переменных окружения
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Формируем DSN из отдельных параметров
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "root")
		dbname := getEnv("DB_NAME", "dictionary")

		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

	// Открываем подключение
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настраиваем пул подключений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
