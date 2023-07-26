package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Chat struct {
	Id int `json:"id"`
}

type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

var telegramApi = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_BOT_TOKEN"))

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		u        Update
		response = events.APIGatewayProxyResponse{StatusCode: 200}
	)

	json.Unmarshal([]byte(request.Body), &u)

	if strings.HasPrefix(u.Message.Text, "/secret") {
		http.PostForm(
			telegramApi,
			url.Values{
				"chat_id": {strconv.Itoa(u.Message.Chat.Id)},
				"text":    {"Ok2"},
			})
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
