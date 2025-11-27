package models

type Content struct {
	Text string `json:"text"`
}

type SMS struct {
	Originator string  `json:"originator"`
	Content    Content `json:"content"`
}

type Message struct {
	Recipient string `json:"recipient"`
	MessageID string `json:"message-id"`
	SMS       SMS    `json:"sms"`
}

type Body struct {
	Messages []Message `json:"messages"`
}
