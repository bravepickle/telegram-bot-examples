package main

import (
	"encoding/json"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

// TODO: make map of structs with handling each type of request
// TODO: use channels for each type of request
// TODO: when logging add prefix of each controller
// TODO: before handling each response spawn different channel to avoid waiting for processing completion

// apiBaseUri base address for Telegram Bots API
const apiBaseUri = `https://api.telegram.org/bot`

// responseTimeoutDefault default timeout for handling telegram requests
const responseTimeoutDefault = 5

/////////////////

type TelegramBotsApiRequestModel struct {
	Path    string                 // URI path of Request
	Timeout int                    // timeout for request
	Api     *TelegramBotsApiStruct // parent for API
}

// Uri builds URI for the request
func (r *TelegramBotsApiRequestModel) Uri() string {
	//fmt.Printf("IN URI: %v\n", r)
	return r.Api.GetBaseUri() + r.Path
	//return `[UNDEFINED]`
}

func (r TelegramBotsApiRequestModel) String() string {
	return r.Path
}

func (r *TelegramBotsApiRequestModel) init() {
	// overreide this method in children
	r.Timeout = responseTimeoutDefault
}

/////////////////

//type MeRequestModel TelegramBotsApiRequestModel
type MeRequestModel struct {
	TelegramBotsApiRequestModel
}

func (r *MeRequestModel) init(api *TelegramBotsApiStruct) {
	r.Path = `/getMe`
	r.Api = api
	r.Timeout = responseTimeoutDefault

	logger.Debug("Initialized Telegram request model: %s", r.Path)
}

/////////////////

type UpdateRequestModel struct {
	TelegramBotsApiRequestModel

	Offset uint32
}

func (r *UpdateRequestModel) init(api *TelegramBotsApiStruct) {
	r.Path = `/getUpdates`
	r.Api = api
	r.Timeout = responseTimeoutDefault

	logger.Debug("Initialized Telegram request model: %s", r.Path)
}

// Uri builds URI for the request
func (r *UpdateRequestModel) Uri() string {
	return r.Api.GetBaseUri() + r.Path + `?timeout=` + strconv.Itoa(r.Timeout) + `&offset=` + strconv.FormatInt(int64(r.Offset), 10)
}

/////////////////

type SendMessageRequestModel struct {
	TelegramBotsApiRequestModel
}

func (r *SendMessageRequestModel) init(api *TelegramBotsApiStruct) {
	r.Path = `/sendMessage`
	r.Api = api
	r.Timeout = responseTimeoutDefault

	logger.Debug("Initialized Telegram request model: %s", r.Path)
}

// Uri builds URI for the request
func (r *SendMessageRequestModel) Uri() string {
	return r.Api.GetBaseUri() + r.Path + `?timeout=` + strconv.Itoa(r.Timeout)
}

/////////////////

type TelegramBotsApiStruct struct {
	BaseUri        string                // base API URI
	RequestManager *RequestManagerStruct // handling requests
	AuthKey        string                // API auth key
	BotInfo        BotProfileTelegramModel
	Sleep          int // interval between polls

	routingMe     MeRequestModel
	routingUpdate UpdateRequestModel
	routingSend   SendMessageRequestModel

	//routing map[string]TelegramBotsApiRequestModel // routing for requests
}

func (r TelegramBotsApiStruct) String() string {
	return r.GetBaseUri()
}

func (r TelegramBotsApiStruct) GetBaseUri() string {
	return r.BaseUri + r.AuthKey
}

func (r *TelegramBotsApiStruct) checkConnection() bool {
	body, ok := r.RequestManager.SendGetRequest(r.routingMe.Uri())
	if !ok {
		return false
	}

	if err := json.Unmarshal(body, &r.BotInfo); err != nil {
		logger.Error("Failed to parse JSON: %s", err)

		return false
	}

	logger.Debug(`Received bot profile: %s`, r.BotInfo)

	return true
}

func (r *TelegramBotsApiStruct) processRequests() {
	logger.Debug(`Starting processing requests...`)

	terminated := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			logger.Info(`Signal "%s" called`, sig)
			terminated = true
		}
	}()

	for {
		r.processUpdates()

		if terminated {
			break
		}

		time.Sleep(time.Duration(r.Sleep) * time.Second)
	}

	logger.Debug(`Finished processing requests.`)
}

func (r *TelegramBotsApiStruct) processUpdates() bool {
	logger.Debug(`Starting polling for updates...`)

	body, ok := r.RequestManager.SendGetRequest(r.routingUpdate.Uri())
	if !ok {
		return false
	}

	var updates UpdateCollectionTelegramModel

	if err := json.Unmarshal(body, &updates); err != nil {
		logger.Error("Error parsing JSON: %s", err)

		return false
	}

	// TODO: edited_message handle, inline_query

	sentOnceSuccessfully := false

	for _, upd := range updates.Result {
		logger.Info(`Handling update `, upd.UpdateId, `message`, upd.Message.MessageId)
		logger.Debug(`> %s`, upd.Message.Text)

		var text = upd.Message.Text
		for _, ent := range upd.Message.Entities {
			if ent.Type == `bot_command` {
				cmd := upd.Message.Text[ent.Offset : ent.Offset+ent.Length]
				logger.Debug(`Is bot command:`, cmd)

				switch cmd {
				case `/start`:
					text = `Hi, ` + upd.Message.From.FirstName + ` ` + upd.Message.From.LastName + `. Thanks for using this bot!`
					//case `/time`:
					//	text = `*Bot time:* ` + time.Now().Format("2006-01-02 15:04:05")
					//case `/code`:
					//	// log.Println(`>>> "`+upd.Message.Text+`"`, ent.Offset+ent.Length+1, ent)
					//	if len(upd.Message.Text) <= ent.Offset+ent.Length+1 {
					//		text = `No input...`
					//	} else {
					//		text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
					//		text = "```\n" + text + "\n```"
					//	}
					//
					//	// text = `*Bot time:* ` + time.Now().Format("2006-01-02 15:04:05")
					//case `/sh`:
					//	query := strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
					//	text = "```\n$ " + query + "\n"
					//
					//	cmd := exec.Command(`/bin/bash`, `-c`, query)
					//	cmd.Env = os.Environ()
					//	out, err := cmd.Output()
					//	if err != nil {
					//		out = []byte(`ERROR: ` + err.Error())
					//	}
					//
					//	text += "\n" + string(out) + "\n```"

				default:
					if len(upd.Message.Text) > ent.Offset+ent.Length {
						text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
					} else {
						text = `Sorry, cannot process your command`
					}
				}
				logger.Debug(`Message changed to:`, text)
			} else if !ent.allowedType() {
				logger.Info(`Warning! Unexpected MessageEntity type:`, ent.Type)
			}
		}

		msg := NewSendMessage(upd.Message.Chat.Id, text /*, upd.Message.MessageId*/)

		payload := url.Values{}
		for name, value := range msg {
			payload.Set(name, value)
		}

		if _, ok := r.RequestManager.SendPostRequest(r.routingSend.Uri(), []byte(payload.Encode())); !ok {
			logger.Error("Failed to send message: %s", payload)
		}

		sentOnceSuccessfully = true

		if upd.UpdateId >= r.routingUpdate.Offset {
			logger.Debug("Was offset %d, will be: %d", r.routingUpdate.Offset, upd.UpdateId+1)
			r.routingUpdate.Offset = upd.UpdateId + 1
		}
	}

	logger.Debug(`Finished polling for updates.`)

	return sentOnceSuccessfully
}

func NewSendMessage(chatId uint32, text string /*, replyToMsgId uint32*/) sendMessageStruct {
	msg := make(sendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`disable_notification`] = `true`
	msg[`disable_web_page_preview`] = `true`

	msg[`chat_id`] = strconv.FormatInt(int64(chatId), 10)
	msg[`text`] = text
	//msg[`reply_to_message_id`] = strconv.FormatInt(int64(replyToMsgId), 10)

	return msg
}

/////////////////

func NewTelegramBotsApi(authKey string, sleep int) *TelegramBotsApiStruct {
	var requestManager RequestManagerStruct

	api := TelegramBotsApiStruct{
		BaseUri:        apiBaseUri,
		RequestManager: &requestManager,
		AuthKey:        authKey,
		Sleep:          sleep,
	}

	api.routingMe.init(&api)
	api.routingUpdate.init(&api)
	api.routingSend.init(&api)

	return &api
}
