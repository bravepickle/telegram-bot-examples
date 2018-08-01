package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

	//commands []BotCommand
	commands []BotCommander

	commandDefault DefaultBotCommandStruct

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

func (r *TelegramBotsApiStruct) processScheduledTasks() {
	logger.Debug(`Starting tasks scheduler...`)

	for {
		//tasks := dbManager.findFutureTasks()
		//if logger.DebugLevel() {
		//	logger.Debug(`Scheduled tasks: %s`, encodeToJson(tasks))
		//}

		//r.processUpdates()
		r.processTaskNotifications()

		//logger.Debug(`Sleep...`)

		time.Sleep(time.Duration(r.Sleep) * time.Second) // TODO: change to another value - once a day should be run or similar
	}
}

func (r *TelegramBotsApiStruct) processTaskNotifications() {
	tasks := dbManager.findFutureTasks()

	if logger.DebugLevel() {
		logger.Debug(`Scheduled tasks: %s`, encodeToJson(tasks))
	}

	if len(tasks) > 0 {
		// TODO: send message and somehow define to which user/chat to send it
		// TODO: notify on expired and prompt on what to do - prolong or cancel with xp loss

		tasksByUser := make(map[int][]TaskDbEntity)
		for _, task := range tasks {
			tasksByUser[task.UserId] = append(tasksByUser[task.UserId], task)
		}

		for userId, userTasks := range tasksByUser {
			msg := "*Tasks TODO:*\n"
			// TODO: order by date or priority, exp
			for _, task := range userTasks {
				msg += fmt.Sprintf("  _%s_: expires at \"%s\", gain \"%d\" XP\n", task.Title, task.DateExpiration, task.Exp)
			}

			// TODO: can we use userId as chat id?
			r.sendMessage(NewSendMessage(uint32(userId), msg, 0))
		}
	}
}
func (r *TelegramBotsApiStruct) processRequests() {
	logger.Debug(`Starting processing requests...`)

	terminated := false

	if !logger.DebugLevel() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for sig := range c {
				logger.Info(`Signal "%s" called. Terminating gracefully...`, sig)
				terminated = true
			}
		}()
	}

	for {
		r.processUpdates()

		if terminated {
			break
		}

		logger.Debug(`Sleep...`)

		time.Sleep(time.Duration(r.Sleep) * time.Second)

		if terminated {
			break
		}
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

	//sentOnceSuccessfully := false

	for _, upd := range updates.Result {
		// TODO: add limit of maximum parallel requests to API in parallel and wait until available
		go r.processSingleUpdate(RunOptionsStruct{Upd: upd})
	}

	logger.Debug(`Finished polling for updates.`)

	//return sentOnceSuccessfully

	return true
}

func (r *TelegramBotsApiStruct) processSingleUpdate(options RunOptionsStruct) {
	//var sendMessage sendMessageStruct

	// TODO: use go channels for each update and read channel for results of sentOnceSuccessfully
	// for parallel computation
	logger.Info(`Handling update ID=%d, Message=%d`, options.Upd.UpdateId, options.Upd.Message.MessageId)
	logger.Debug(`> %s`, options.Upd.Message.Text)

	//var text = upd.Message.Text
	var hasProcessed = true

	if len(options.Upd.Message.Entities) > 0 {
		logger.Debug(`Processing message entities...`)
		for _, ent := range options.Upd.Message.Entities {
			options.Ent = ent
			sendMessage, found := r.processMessageEntity(options)

			if found {
				hasProcessed = true

				r.sendMessage(sendMessage)
			}
		}
	} else {
		for _, botCommand := range r.commands {
			if botCommand.IsRunning(options) {
				//// TODO: reset all running commands for user and reinit current one if newly called??
				//found = true
				sendMessage, err := botCommand.Run(options)
				if err != nil {
					logger.Fatal(`Command run "%s" failed: %s`, botCommand.GetName(), err)
				}

				hasProcessed = true

				r.sendMessage(sendMessage)

				//logger.Fatal(`Found running options for command %s`, botCommand.GetName())
				//
				//break
			}
		}
	}

	// todo: update index even if no text messages can be processed

	// TODO: always single send-message for message+message entities?

	if !hasProcessed {
		logger.Debug(`No new messages...`)
	} else if options.Upd.UpdateId >= r.routingUpdate.Offset {
		logger.Debug("Was offset %d, will be: %d", r.routingUpdate.Offset, options.Upd.UpdateId+1)
		r.routingUpdate.Offset = options.Upd.UpdateId + 1
	}
}

func (r *TelegramBotsApiStruct) sendMessage(sendMessage sendMessageStruct) {
	if logger.DebugLevel() {
		logger.Debug(`Message to send: %s`, encodeToJson(sendMessage))
	}

	if _, ok := r.RequestManager.SendPostJsonRequest(r.routingSend.Uri(), sendMessage); !ok {
		logger.Error("Failed to send message: %s", encodeToJson(sendMessage))
	}
}

func (r *TelegramBotsApiStruct) processMessageEntity(runOptions RunOptionsStruct) (sendMessage sendMessageStruct, found bool) {
	var err error
	found = false

	if runOptions.Ent.Type == `bot_command` {
		cmd := runOptions.Upd.Message.Text[runOptions.Ent.Offset : runOptions.Ent.Offset+runOptions.Ent.Length]
		logger.Debug(`Is bot command: %s`, cmd)

		for _, botCommand := range r.commands {
			if botCommand.IsRunning(runOptions) || cmd == botCommand.GetName() {
				// TODO: reset all running commands for user and reinit current one if newly called??
				found = true
				sendMessage, err = botCommand.Run(runOptions)
				if err != nil {
					logger.Fatal(`Command run "%s" failed: %s`, botCommand.GetName(), err)
				}

				break
			}
		}

		if !found {
			logger.Debug(`Running default command`)
			sendMessage, err = r.commandDefault.Run(runOptions)
			if err != nil {
				logger.Fatal(`Command run "%s" failed: %s`, r.commandDefault.GetName(), err)
			}
		}
	} else if !runOptions.Ent.allowedType() {
		logger.Info(`Warning! Unexpected MessageEntity type: %s`, runOptions.Ent.Type)
	} else {
		logger.Info(`Failed to handle message entity type: %s`, runOptions.Ent.Type)
	}

	return sendMessage, found
}

func NewSendMessage(chatId uint32, text string, replyToMsgId uint32) sendMessageStruct {
	msg := make(sendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`disable_notification`] = `true`
	msg[`disable_web_page_preview`] = `true`

	msg[`chat_id`] = strconv.FormatInt(int64(chatId), 10)
	msg[`text`] = text

	if replyToMsgId != 0 {
		msg[`reply_to_message_id`] = replyToMsgId
	}

	return msg
}

type SendMessageOptionsStruct struct {
	ChatId       int
	Text         string
	ReplyToMsgId int
}

// NewSendMessageWithOptions generates send message according to input options
func NewSendMessageWithOptions(options SendMessageOptionsStruct) sendMessageStruct {
	msg := make(sendMessageStruct)

	msg[`parse_mode`] = `Markdown`
	msg[`disable_notification`] = `true`
	msg[`disable_web_page_preview`] = `true`

	msg[`chat_id`] = options.ChatId
	msg[`text`] = options.Text

	if options.ReplyToMsgId != 0 {
		msg[`reply_to_message_id`] = options.ReplyToMsgId
	}

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

	api.commands = append(api.commands, NewStartBotCommand(), NewAddTaskBotCommand())

	return &api
}
