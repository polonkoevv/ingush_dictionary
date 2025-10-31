package application

import (
	"context"
	"database/sql"
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

	dicts, err := s.dictRepo.GetAll(ctx)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return dicts, nil
}
