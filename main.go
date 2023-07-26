package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/valyala/fastjson"
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
		p        fastjson.Parser
		response events.APIGatewayProxyResponse
	)

	v, err := p.Parse(request.Body)
	if err != nil {
		return response, err
	}

	fmt.Printf("Chat ID2: %d\n", v.GetUint64("message", "chat", "id"))

	/*
		response, err := http.PostForm(
			telegramApi,
			url.Values{
				"chat_id": {strconv.Itoa(update.Message.Chat.Id)},
				"text":    {"Ok!"},
			})

		fmt.Printf("Response %v\n", response)
		body, err := io.ReadAll(response.Body)

		fmt.Printf("Response Body %s\n", string(body[:]))
		fmt.Printf("Error %v\n", err)
	*/

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
