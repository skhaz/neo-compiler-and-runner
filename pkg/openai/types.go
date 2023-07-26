package openai

import "net/http"

type RoleType string

type ModelType string

type Message struct {
	Role    RoleType `json:"role,omitempty"`
	Content string   `json:"content"`
}

type Request struct {
	Model    ModelType  `json:"model"`
	Messages []*Message `json:"messages"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

type Response struct {
	Choices []*Choice `json:"choices"`
	Error   *Error    `json:"error,omitempty"`
}

type Client struct {
	client *http.Client
	apiKey string
	url    string
}
