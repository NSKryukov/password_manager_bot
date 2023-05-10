package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Text      string `json:"text"`
	Chat      struct {
		ChatId int `json:"id"`
	}
	UserInfo struct {
		Username string `json:"username"`
	} `json:"from"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type MessageToSend struct {
	Text   string `json:"text"`
	ChatId int    `json:"chat_id"`
}

type MessageToDelete struct {
	ChatID    int `json:"chat_id"`
	MessageID int `json:"message_id"`
}

type ReceivedMessage struct {
	Result Message `json:"result"`
}
