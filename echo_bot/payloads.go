package main

import "strconv"

type messageStruct struct {
	UpdateId uint32 `json:"update_id"`
	Message struct {
		MessageId uint32 `json:"message_id"`
		From struct {
			Id           uint32
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			//Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		}
		Chat struct {
			Id        uint32
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Type      string
		}
		Date uint32
		Text string
	}
}

type updatesPayloadStruct struct {
	Ok bool
	Result []messageStruct
}



type sendMessageStruct map[string]string
//type sendMessageStruct struct {
//	ChatId uint32 `json:"chat_id"`
//	Text string `json:"text"`
//	ParseMode string `json:"parse_mode"`
//	DisableWebPagePreview bool `json:"disable_web_page_preview"`
//	DisableNotification bool `json:"disable_notification"`
//	ReplyToMessageId uint32 `json:"reply_to_message_id"`
//}

func NewSendMessage(chatId uint32, text string, replyToMsgId uint32) sendMessageStruct {
	msg := make(sendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`disable_notification`] = `true`
	msg[`disable_web_page_preview`] = `true`

	msg[`chat_id`] = strconv.FormatInt(int64(chatId), 10)
	msg[`text`] = text
	msg[`reply_to_message_id`] = strconv.FormatInt(int64(replyToMsgId), 10)
	//msg[`timeout`] = strconv.Itoa(5)

	return msg
}

/**
{
  "ok": true,
  "result": [
    {
      "update_id": 574275051,
      "message": {
        "message_id": 8,
        "from": {
          "id": 209139256,
          "first_name": "Victor",
          "last_name": "K",
          "language_code": "en-UA"
        },
        "chat": {
          "id": 209139256,
          "first_name": "Victor",
          "last_name": "K",
          "type": "private"
        },
        "date": 1498684667,
        "text": "hello me"
      }
    },
    {
      "update_id": 574275052,
      "message": {
        "message_id": 9,
        "from": {
          "id": 209139256,
          "first_name": "Victor",
          "last_name": "K",
          "language_code": "en-UA"
        },
        "chat": {
          "id": 209139256,
          "first_name": "Victor",
          "last_name": "K",
          "type": "private"
        },
        "date": 1498684723,
        "text": "hi, me"
      }
    }
  ]
}
 */
