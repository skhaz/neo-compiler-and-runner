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
	prefix = "/run"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		h      = openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		update = telegram.Parse(request.Body)
	)

	if strings.HasPrefix(update.Message.Text, prefix) {
		request := &openai.Request{
			Model: openai.ModelGpt35Turbo,
			Messages: []*openai.Message{
				{Role: openai.RoleSystem, Content: "Keep strictly to say only the output, without any comment, what is the result of the following code?"},
				{Role: openai.RoleUser, Content: strings.Trim(update.Message.Text, prefix)},
			}}

		response, err := h.Do(request)
		if err != nil || response.Error != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		telegram.Reply(os.Getenv("TELEGRAM_BOT_TOKEN"), update.Message.Chat.Id, response.Choices[0].Message.Content)
	}

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
