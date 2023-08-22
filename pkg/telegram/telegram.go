package telegram

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func Parse(body string) Update {
	var update Update
	json.Unmarshal([]byte(body), &update)

	return update
}

func Reply(token string, id int, text string) {
	var endpoint string = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	http.PostForm(endpoint, url.Values{"chat_id": {strconv.Itoa(id)}, "text": {text}})
}
