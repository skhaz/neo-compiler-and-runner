package main

import (
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"skhaz.dev/compliquer/pkg/openai"
	"skhaz.dev/compliquer/pkg/telegram"
)

const (
	prefix = "/secret"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		h        = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		update   = telegram.Parse(request.Body)
		response = events.APIGatewayProxyResponse{StatusCode: 200}
	)

	if strings.HasPrefix(update.Message.Text, prefix) {
		request := &openai.Request{
			Model: openai.ModelGpt35Turbo,
			Messages: []*openai.Message{
				{Role: openai.RoleSystem, Content: "You are a compiler assistant who compiles or interpret code."},
				{Role: openai.RoleSystem, Content: "What is the output of the following code?"},
				{Role: openai.RoleUser, Content: strings.Trim(update.Message.Text, prefix)},
			}}

		r, err := h.Do(request)
		if err != nil {
			return response, err
		}

		if r.Error != nil {
			return response, err
		}

		telegram.Reply(os.Getenv("TELEGRAM_BOT_TOKEN"), update.Message.Chat.Id, r.Choices[0].Message.Content)
	}

	return response, nil
}

func main() {
	lambda.Start(Handler)
}
