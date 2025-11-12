package application

import (
	"test/internal/infrustructure/postgres"
)

type Service struct {
	UserSrv UserService
	DictSrv DictService
	WordSrv WordService
}

func CreateServices(rep postgres.Repository) (*Service, error) {
	userSrv := NewUserService(rep.UserRep)
	wordSrv := NewWordService(rep.WordRep)
	dictSrv := NewDictService(rep.DictRep)

	return &Service{
		UserSrv: *userSrv,
		WordSrv: *wordSrv,
		DictSrv: *dictSrv,
	}, nil
}
