package inmemory

type InMemoryStorage struct {
	language map[int64]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		language: make(map[int64]string),
	}
}

func (s *InMemoryStorage) ChangeLanguage(chatID int64) (string, error) {
	if s.language[chatID] == "rus" {
		s.language[chatID] = "ing"
	} else {
		s.language[chatID] = "rus"
	}
	return s.language[chatID], nil
}

func (s *InMemoryStorage) GetLanguage(chatID int64) (string, error) {
	language, ok := s.language[chatID]
	if !ok {
		s.language[chatID] = "ing"
		return "ing", nil
	}
	return language, nil
}
