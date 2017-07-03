package main

import "strconv"

type userStruct struct {
	Id           uint32
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type chatStruct struct {
	Id        uint32
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Type      string `json:"type"`
	Title     string `json:"title"`
}

type messageEntityStruct struct {
	// Type of the entity. Can be mention (@username), hashtag, bot_command, url, email, bold (bold text), italic (italic text), code (monowidth string), pre (monowidth block), text_link (for clickable text URLs), text_mention (for users without usernames)
	Type string
	// Offset in UTF-16 code units to the start of the entity
	Offset int
	// Length of the entity in UTF-16 code units
	Length int
	//Optional. For “text_link” only, url that will be opened after user taps on the text
	Url string
	//Optional. For “text_mention” only, the mentioned user
	User userStruct
}

func (ms *messageEntityStruct) allowedType() bool {
	allowed := []string{`bot_command`, `url`, `email`, `code`}

	for _, cmd := range allowed {
		if ms.Type == cmd {
			return true
		}
	}

	return false
}

type messageStruct struct {
	MessageId uint32 `json:"message_id"`
	From      userStruct
	Chat      chatStruct
	Date      uint32 `json:"date"`
	Text      string `json:"text"`
	Entities  []messageEntityStruct `json:"entities"`
}

type updateStruct struct {
	UpdateId uint32 `json:"update_id"`
	Message  messageStruct
}

type updatesPayloadStruct struct {
	Ok     bool
	Result []updateStruct
}

type sendMessageStruct map[string]string

func NewSendMessage(chatId uint32, text string/*, replyToMsgId uint32*/) sendMessageStruct {
	msg := make(sendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`disable_notification`] = `true`
	msg[`disable_web_page_preview`] = `true`

	msg[`chat_id`] = strconv.FormatInt(int64(chatId), 10)
	msg[`text`] = text
	//msg[`reply_to_message_id`] = strconv.FormatInt(int64(replyToMsgId), 10)

	return msg
}
