package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

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

func Handler(update Update) (Response, error) {
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

	return Response{Ok: true}, nil
	/*
		var buf bytes.Buffer

		body, err := json.Marshal(map[string]interface{}{
			"message": "Go Serverless v1.0! Your function executed successfully!",
		})
		if err != nil {
			return Response{StatusCode: 404}, err
		}
		json.HTMLEscape(&buf, body)

		resp := Response{
			StatusCode:      200,
			IsBase64Encoded: false,
			Body:            buf.String(),
			Headers: map[string]string{
				"Content-Type":           "application/json",
				"X-MyCompany-Func-Reply": "hello-handler",
			},
		}

		return resp, nil
	*/
}

func main() {
	lambda.Start(Handler)
}
