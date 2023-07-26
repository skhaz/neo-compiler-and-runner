package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const URL = "https://api.openai.com/v1/chat/completions"

const ModelGpt35Turbo = "gpt-3.5-turbo"

const (
	RoleUser      RoleType = "user"
	RoleAssistant RoleType = "assistant"
	RoleSystem    RoleType = "system"
)

func NewClient(apiKey string) *Client {
	return &Client{
		client: &http.Client{},
		apiKey: apiKey,
		url:    URL,
	}
}

func (c *Client) Do(r *Request) (*Response, error) {
	var buffer = bytes.NewBuffer(nil)
	if err := json.NewEncoder(buffer).Encode(r); err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}

	request, err := http.NewRequest(http.MethodPost, c.url, buffer)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http do: %w", err)
	}

	defer func() {
		io.Copy(io.Discard, response.Body)
		response.Body.Close()
	}()

	var result Response
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	return &result, nil
}
