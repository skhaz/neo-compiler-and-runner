package main

import (
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

func Handler(update Update) (Response, error) {

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
