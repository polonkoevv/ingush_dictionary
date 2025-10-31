package application

import (
	"context"
	"errors"
	"test/internal/domain/word"
)

type WordService struct {
	wordRepo word.Repository
}

func NewWordService(wordRepo word.Repository) *WordService {
	return &WordService{
		wordRepo: wordRepo,
	}
}

func (s *WordService) GetTranslation(ctx context.Context, query string, lang_from string) ([]*word.Word, error) {
	switch lang_from {
	case "ing":
		return s.wordRepo.GetIng(ctx, query)
	case "rus":
		return s.wordRepo.GetRus(ctx, query)
	default:
		return nil, errors.New("got invalid language")
	}
}

func (s *WordService) GetTranslationFiltered(ctx context.Context, query string, lang_from string, tg_user_id int64) ([]*word.Word, error) {
	switch lang_from {
	case "ing":
		return s.wordRepo.GetIngFiltered(ctx, query, tg_user_id)
	case "rus":
		return s.wordRepo.GetRusFiltered(ctx, query, tg_user_id)
	default:
		return nil, errors.New("got invalid language")
	}
}
