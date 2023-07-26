package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"skhaz.dev/compliquer/pkg/telegram"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		// openai   = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		update   = telegram.Parse(request.Body)
		response = events.APIGatewayProxyResponse{StatusCode: 200}
	)

	if strings.HasPrefix(update.Message.Text, "/secret") {
		code := strings.Trim(update.Message.Text, "/secret ")

		fmt.Printf("Code %s\n", code)

		telegram.Reply(os.Getenv("TELEGRAM_API_KEY"), update.Message.Chat.Id, code)
		/*
			http.PostForm(
				telegramApi,
				url.Values{
					"chat_id": {strconv.Itoa(u.Message.Chat.Id)},
					"text":    {"Ok2"},
				})
		*/
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
