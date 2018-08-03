package main

import (
	"encoding/json"
	"strconv"
)

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

type ArrayEncoder interface {
	ToArray() map[string]interface{}
}

type SendMessageOptionsStruct struct {
	ChatId                int
	Text                  string
	ReplyToMsgId          int
	DisableWebPagePreview bool
	DisableNotification   bool
	ReplyMarkup           ArrayEncoder
}

type InlineKeyboardButtonTelegramModel struct {
	Text                         string
	Url                          string
	CallbackData                 string
	SwitchInlineQuery            string
	SwitchInlineQueryCurrentChat string
	Pay                          bool
}

func (k InlineKeyboardButtonTelegramModel) ToArray() map[string]interface{} {
	v := make(map[string]interface{})

	v[`text`] = k.Text

	if k.Url != `` {
		v[`url`] = k.Url
	}

	if k.CallbackData != `` {
		v[`callback_data`] = k.CallbackData
	}

	if k.SwitchInlineQuery != `` {
		v[`switch_inline_query`] = k.SwitchInlineQuery
	}

	if k.SwitchInlineQueryCurrentChat != `` {
		v[`switch_inline_query_current_chat`] = k.SwitchInlineQueryCurrentChat
	}

	if k.Pay {
		v[`switch_inline_query_current_chat`] = `true`
	}

	return v
}

type InlineKeyboardMarkupTelegramModel struct {
	InlineKeyboard []InlineKeyboardButtonTelegramModel
}

func (k InlineKeyboardMarkupTelegramModel) ToArray() map[string]interface{} {
	v := make(map[string]interface{})

	//type Keyboards []map[string]interface{}

	//var inlineKeyboards []map[string]interface{}
	var inlineKeyboards []map[string]interface{}
	var keyboardsCollection [][]map[string]interface{}

	//v[`inline_keyboard`] = make(map[string]interface{})

	//v[`inline_keyboard`] = k.InlineKeyboard.ToArray()

	for _, kb := range k.InlineKeyboard {
		//v[`inline_keyboard`] = append(v[`inline_keyboard`], kb.ToArray())
		//v[`inline_keyboard`] = append(v[`inline_keyboard`], kb.ToArray())
		inlineKeyboards = append(inlineKeyboards, kb.ToArray())
	}

	keyboardsCollection = append(keyboardsCollection, inlineKeyboards)
	v[`inline_keyboard`] = keyboardsCollection

	//v[`inline_keyboard`] = k.InlineKeyboard.ToArray()

	return v
}

func NewSendMessage(chatId uint32, text string, replyToMsgId uint32) SendMessageStruct {
	msg := make(SendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`disable_notification`] = true
	msg[`disable_web_page_preview`] = true

	msg[`chat_id`] = strconv.FormatInt(int64(chatId), 10)
	msg[`text`] = text

	if replyToMsgId != 0 {
		msg[`reply_to_message_id`] = replyToMsgId
	}

	return msg
}

// NewSendMessageWithOptions generates send message according to input options
func NewSendMessageWithOptions(options SendMessageOptionsStruct) SendMessageStruct {
	msg := make(SendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`chat_id`] = options.ChatId
	msg[`text`] = options.Text

	if options.DisableNotification {
		msg[`disable_notification`] = true
	}

	if options.DisableWebPagePreview {
		msg[`disable_web_page_preview`] = true
	}

	if options.ReplyToMsgId != 0 {
		msg[`reply_to_message_id`] = options.ReplyToMsgId
	}

	if options.ReplyMarkup != nil {
		msg[`reply_markup`] = options.ReplyMarkup.ToArray()
	}

	return msg
}
