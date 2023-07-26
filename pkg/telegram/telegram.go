package telegram

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func Parse(body string) Update {
	var update Update
	json.Unmarshal([]byte(body), &update)

	return update
}

func Reply(token string, id int, text string) error {
	endpoint := "https://api.telegram.org/bot" + token + "/sendMessage"

	_, err := http.PostForm(
		endpoint,
		url.Values{
			"chat_id": {strconv.Itoa(id)},
			"text":    {text},
		})

	return err
}
