package handler

import (
	"skhaz.dev/compliquer/pkg/openai"
)

type Handler struct {
	OpenAI *openai.Client
}

type Error struct {
	Error string `json:"error"`
}

type Response struct {
	Result string `json:"result"`
}
