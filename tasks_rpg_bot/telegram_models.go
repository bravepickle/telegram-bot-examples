package main

import "encoding/json"

// Telegram Bots API models

type UserTelegramModel struct {
	Id           uint32
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type ChatTelegramModel struct {
	Id        uint32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
	Title     string `json:"title"`
}

type MessageEntityTelegramModel struct {
	// Type of the entity. Can be mention (@username), hashtag, bot_command, url, email, bold (bold text), italic (italic text), code (monowidth string), pre (monowidth block), text_link (for clickable text URLs), text_mention (for users without usernames)
	Type string
	// Offset in UTF-16 code units to the start of the entity
	Offset int
	// Length of the entity in UTF-16 code units
	Length int
	//Optional. For “text_link” only, url that will be opened after user taps on the text
	Url string
	//Optional. For “text_mention” only, the mentioned user
	User UserTelegramModel
}

func (ms *MessageEntityTelegramModel) allowedType() bool {
	allowed := []string{`bot_command`, `url`, `email`, `code`}

	for _, cmd := range allowed {
		if ms.Type == cmd {
			return true
		}
	}

	return false
}

type BotProfileTelegramModel struct {
	Ok     bool
	Result struct {
		Id        uint32
		FirstName string `json:"first_name"`
		Username  string
	}
}

func (p BotProfileTelegramModel) String() string {
	if text, err := json.Marshal(p); err == nil {
		return string(text)
	} else {
		logger.Error(`Error ocurred when encoding to JSON bot profile object: %s`, err)

		return `[UNDEFINED]`
	}
}

type MessageTreeTelegramModel struct {
	MessageId uint32 `json:"message_id"`
	From      UserTelegramModel
	Chat      ChatTelegramModel
	Date      uint32                       `json:"date"`
	Text      string                       `json:"text"`
	Entities  []MessageEntityTelegramModel `json:"entities"`
}

type UpdateTelegramModel struct {
	UpdateId uint32 `json:"update_id"`
	Message  MessageTreeTelegramModel
}

type UpdateCollectionTelegramModel struct {
	Ok     bool
	Result []UpdateTelegramModel
}

type SendMessageStruct map[string]interface{}
