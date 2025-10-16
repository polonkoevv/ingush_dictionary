package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"test/internal/models"
	"test/internal/storage"
)

const BASE_URL = "https://doshlorg.ru/api/words?search=%s&language=%s&limit=100"
const BASE_URL_RUS = "https://doshlorg.ru/api/words?search=%s&language=rus&limit=100"
const BASE_URL_ING = "https://doshlorg.ru/api/words?search=%s&language=ing&limit=100"

type LanguageService interface {
	GetLanguage(chatID int64) (string, error)
	ChangeLanguage(chatID int64) (string, error)
	IngToRus(word string) (string, int, error)
	RusToIng(word string) (string, int, error)
}

type languageService struct {
	storage storage.Storage
}

func PrepareWord(word string) string {
	word = strings.ReplaceAll(word, "1", "i")
	word = strings.ToLower(word)
	word = strings.TrimSpace(word)
	return word
}

func NewLanguageService(storage storage.Storage) LanguageService {
	return &languageService{storage: storage}
}

func (s *languageService) GetLanguage(chatID int64) (string, error) {
	return s.storage.GetLanguage(chatID)
}

func (s *languageService) ChangeLanguage(chatID int64) (string, error) {
	return s.storage.ChangeLanguage(chatID)
}

func (s *languageService) IngToRus(word string) (string, int, error) {

	word = PrepareWord(word)

	res, err := http.Get(fmt.Sprintf(BASE_URL_ING, word))
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	var wordResponse models.WordResponse
	err = json.Unmarshal(body, &wordResponse)
	if err != nil {
		return "", 0, err
	}

	if len(wordResponse.Data) == 0 {
		return "", 0, nil
	}

	rs := ""

	for _, wrd := range wordResponse.Data {
		rs = rs + fmt.Sprintf("%s\n", wrd.Word)
		for _, wrd2 := range wrd.Translates {
			rs = rs + fmt.Sprintf("--- \t%s\n", wrd2.Word)
		}
	}

	return rs, wordResponse.Total, nil
}

func (s *languageService) RusToIng(word string) (string, int, error) {
	res, err := http.Get(fmt.Sprintf(BASE_URL_RUS, word))
	if err != nil {
		return "", 0, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", 0, err
	}

	var wordResponse models.WordResponse
	err = json.Unmarshal(body, &wordResponse)
	if err != nil {
		return "", 0, err
	}

	if len(wordResponse.Data) == 0 {
		return "", 0, nil
	}

	rs := ""

	for _, wrd := range wordResponse.Data {
		rs = rs + fmt.Sprintf("%s\n", wrd.Word)
		for _, wrd2 := range wrd.Words {
			rs = rs + fmt.Sprintf("--- \t%s\n", wrd2.Word)
		}
	}
	return rs, wordResponse.Total, nil
}
