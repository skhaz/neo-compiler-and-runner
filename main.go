package main

import (
	"encoding/json"
	"fmt"
	"os"

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

type Response struct {
	Ok bool `json:"ok"`
}

var telegramApi = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_BOT_TOKEN"))

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		u        Update
		response events.APIGatewayProxyResponse
	)

	err := json.Unmarshal([]byte(request.Body), &u)
	if err != nil {
		return response, err
	}

	fmt.Printf("Chat ID3: %d\n", u.Message.Chat.Id)

	/*
		x, err := http.PostForm(
			telegramApi,
			url.Values{
				"chat_id": {strconv.Itoa(u.Message.Chat.Id)},
				"text":    {"Ok!"},
			})

		fmt.Printf("Response %v\n", response)
		body, err := io.ReadAll(x.Body)

		fmt.Printf("Response Body %s\n", string(body[:]))
		fmt.Printf("Error %v\n", err)
	*/
	return response, nil
}

func main() {
	lambda.Start(Handler)
}
