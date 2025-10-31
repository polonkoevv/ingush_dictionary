package application

import (
	"test/internal/domain/dictionary"
	"test/internal/domain/user"
	"test/internal/domain/word"
)

func CreateServices(userRep *user.Repository, wordRep *word.Repository, dictRep *dictionary.Repository) (*UserService, *WordService, *DictService, error) {
	userSrv := NewUserService(*userRep)
	wordSrv := NewWordService(*wordRep)
	dictSrv := NewDictService(*dictRep)

	return userSrv, wordSrv, dictSrv, nil
}
