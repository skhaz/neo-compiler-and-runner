package telegram

import "net/http"

type Client struct {
	client *http.Client
	url    string
}

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
