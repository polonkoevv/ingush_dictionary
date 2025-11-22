package application

import (
	"context"
	"database/sql"
	"log/slog"
	"test/internal/domain/dictionary"
)

type DictService struct {
	dictRepo dictionary.Repository
}

func NewDictService(dictRepo dictionary.Repository) *DictService {
	return &DictService{
		dictRepo: dictRepo,
	}
}

func (s *DictService) GetAllDictionaries(ctx context.Context) ([]*dictionary.Dictionary, error) {
	slog.Debug("getting all dictionaries",
		slog.String("component", "dict_service"),
		slog.String("op", "GetAllDictionaries"),
	)

	dicts, err := s.dictRepo.GetAll(ctx)

	if err == sql.ErrNoRows {
		slog.Debug("no dictionaries found",
			slog.String("component", "dict_service"),
			slog.String("op", "GetAllDictionaries"),
		)
		return nil, nil
	}

	if err != nil {
		slog.Error("failed to get all dictionaries",
			slog.String("component", "dict_service"),
			slog.String("op", "GetAllDictionaries"),
			slog.Any("error", err),
		)
		return nil, err
	}

	slog.Debug("dictionaries retrieved",
		slog.String("component", "dict_service"),
		slog.String("op", "GetAllDictionaries"),
		slog.Int("dicts_count", len(dicts)),
	)

	return dicts, nil
}
