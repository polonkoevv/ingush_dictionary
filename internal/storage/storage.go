package storage

type Storage interface {
	ChangeLanguage(chatID int64) (string, error)
	GetLanguage(chatID int64) (string, error)
}
