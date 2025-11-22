package application

import (
	"context"
	"errors"
	"log/slog"
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

func (s *WordService) GetTranslation(ctx context.Context, query string, lang_from string, page_number int) ([]*word.Word, int, error) {
	slog.Debug("translation request",
		slog.String("component", "word_service"),
		slog.String("op", "GetTranslation"),
		slog.String("query", query),
		slog.String("lang_from", lang_from),
		slog.Int("page_number", page_number),
	)

	var words []*word.Word
	var totalCount int
	var err error

	switch lang_from {
	case "ing":
		words, totalCount, err = s.wordRepo.GetIng(ctx, query, page_number)
	case "rus":
		words, totalCount, err = s.wordRepo.GetRus(ctx, query, page_number)
	default:
		err = errors.New("got invalid language")
		slog.Error("invalid language in translation request",
			slog.String("component", "word_service"),
			slog.String("op", "GetTranslation"),
			slog.String("query", query),
			slog.String("lang_from", lang_from),
			slog.Int("page_number", page_number),
		)
		return nil, 0, err
	}

	if err != nil {
		slog.Error("failed to get translation",
			slog.String("component", "word_service"),
			slog.String("op", "GetTranslation"),
			slog.String("query", query),
			slog.String("lang_from", lang_from),
			slog.Int("page_number", page_number),
			slog.Any("error", err),
		)
		return nil, 0, err
	}

	slog.Debug("translation retrieved",
		slog.String("component", "word_service"),
		slog.String("op", "GetTranslation"),
		slog.String("query", query),
		slog.String("lang_from", lang_from),
		slog.Int("page_number", page_number),
		slog.Int("words_count", len(words)),
		slog.Int("total_count", totalCount),
	)

	return words, totalCount, nil
}

func (s *WordService) GetTranslationFiltered(ctx context.Context, query string, lang_from string, tg_user_id int64, page_number int) ([]*word.Word, int, error) {
	slog.Debug("filtered translation request",
		slog.String("component", "word_service"),
		slog.String("op", "GetTranslationFiltered"),
		slog.String("query", query),
		slog.String("lang_from", lang_from),
		slog.Int64("tg_user_id", tg_user_id),
		slog.Int("page_number", page_number),
	)

	var words []*word.Word
	var totalCount int
	var err error

	switch lang_from {
	case "ing":
		words, totalCount, err = s.wordRepo.GetIngFiltered(ctx, query, tg_user_id, page_number)
	case "rus":
		words, totalCount, err = s.wordRepo.GetRusFiltered(ctx, query, tg_user_id, page_number)
	default:
		err = errors.New("got invalid language")
		slog.Error("invalid language in filtered translation request",
			slog.String("component", "word_service"),
			slog.String("op", "GetTranslationFiltered"),
			slog.String("query", query),
			slog.String("lang_from", lang_from),
			slog.Int64("tg_user_id", tg_user_id),
			slog.Int("page_number", page_number),
		)
		return nil, 0, err
	}

	if err != nil {
		slog.Error("failed to get filtered translation",
			slog.String("component", "word_service"),
			slog.String("op", "GetTranslationFiltered"),
			slog.String("query", query),
			slog.String("lang_from", lang_from),
			slog.Int64("tg_user_id", tg_user_id),
			slog.Int("page_number", page_number),
			slog.Any("error", err),
		)
		return nil, 0, err
	}

	slog.Debug("filtered translation retrieved",
		slog.String("component", "word_service"),
		slog.String("op", "GetTranslationFiltered"),
		slog.String("query", query),
		slog.String("lang_from", lang_from),
		slog.Int64("tg_user_id", tg_user_id),
		slog.Int("page_number", page_number),
		slog.Int("words_count", len(words)),
		slog.Int("total_count", totalCount),
	)

	return words, totalCount, nil
}
