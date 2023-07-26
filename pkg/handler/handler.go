package handler

import (
	"skhaz.dev/compliquer/pkg/openai"
)

func NewHandler(openai *openai.Client) *Handler {
	return &Handler{
		OpenAI: openai,
	}
}
