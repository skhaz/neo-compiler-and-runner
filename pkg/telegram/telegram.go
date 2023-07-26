package telegram

import (
	"encoding/json"
	"fmt"
	"io"
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
	var endpoint string = "https://api.telegram.org/bot" + token + "/sendMessage"

	response, err := http.PostForm(
		endpoint,
		url.Values{
			"chat_id": {strconv.Itoa(id)},
			"text":    {text},
		})

	r, err := io.ReadAll(response.Body)
	fmt.Printf("Response %s\n", string(r[:]))

	return err
}
