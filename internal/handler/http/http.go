package http

import (
	"test/internal/service"
)

type HttpHandler struct {
	service service.LanguageService
}

func NewHttpHandler(service service.LanguageService) *HttpHandler {
	return &HttpHandler{service: service}
}

func (h *HttpHandler) GetWord(word string, language string) (string, int, error) {
	if language == "rus" {
		return h.service.RusToIng(word)
	}

	return h.service.IngToRus(word)
}
