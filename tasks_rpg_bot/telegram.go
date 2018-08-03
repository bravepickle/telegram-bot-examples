package main

import (
	"encoding/json"
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

/////////////////

type TelegramBotsApiRequestModel struct {
	Path    string                 // URI path of Request
	Timeout int                    // timeout for request
	Api     *TelegramBotsApiStruct // parent for API
}

// Uri builds URI for the request
func (r *TelegramBotsApiRequestModel) Uri() string {
	return r.Api.GetBaseUri() + r.Path
}

func (r TelegramBotsApiRequestModel) String() string {
	return r.Path
}

func (r *TelegramBotsApiRequestModel) Init(api *TelegramBotsApiStruct) {
	// override this method in children
	r.Timeout = appConfig.GetApiTimeout()
	r.Api = api
}

/////////////////

//type MeRequestModel TelegramBotsApiRequestModel
type MeRequestModel struct {
	TelegramBotsApiRequestModel
}

func (r *MeRequestModel) Init(api *TelegramBotsApiStruct) {
	r.TelegramBotsApiRequestModel.Init(api)
	r.Path = `/getMe`

	logger.Debug("Initialized Telegram request model: %s", r.Path)
}

/////////////////

type UpdateRequestModel struct {
	TelegramBotsApiRequestModel

	Offset uint32
}

func (r *UpdateRequestModel) Init(api *TelegramBotsApiStruct) {
	r.TelegramBotsApiRequestModel.Init(api)
	r.Path = `/getUpdates`

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

func (r *SendMessageRequestModel) Init(api *TelegramBotsApiStruct) {
	r.TelegramBotsApiRequestModel.Init(api)
	r.Path = `/sendMessage`

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
		r.processTaskNotifications()
		time.Sleep(time.Duration(appConfig.GetApiRemindInterval()) * time.Hour)
	}
}

func (r *TelegramBotsApiStruct) processTaskNotifications() {
	tasks := dbManager.findFutureTasks()

	if logger.DebugLevel() {
		logger.Debug(`Scheduled tasks: %s`, encodeToJson(tasks))
	}

	if len(tasks) > 0 {
		// TODO: notify on expired and prompt on what to do - prolong or cancel with xp loss?

		tasksByUser := make(map[int][]TaskDbEntity)
		for _, task := range tasks {
			tasksByUser[task.UserId] = append(tasksByUser[task.UserId], task)
		}

		for userId, userTasks := range tasksByUser {
			msg := "*" + usrMsg.T(`remind.task.header`) + "*\n"
			// TODO: order by date or priority, XP
			for k, task := range userTasks {
				msg += usrMsg.T(`remind.task.line`, k+1, task.Title, task.DateExpiration, task.Exp) + "\n"
			}

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

	// TODO: process callback query

	if logger.DebugLevel() {
		logger.Debug(`>>>>>>>> Parsed updates: %s`, encodeToJson(updates))
	}

	// TODO: edited_message handle, inline_query

	for _, upd := range updates.Result {
		// TODO: add limit of maximum parallel requests to API in parallel and wait until available
		go r.processSingleUpdate(RunOptionsStruct{Upd: upd})
	}

	logger.Debug(`Finished polling for updates.`)

	return true
}

func (r *TelegramBotsApiStruct) processSingleUpdate(options RunOptionsStruct) {
	// TODO: use go channels for each update and read channel for results of sentOnceSuccessfully - for parallel computation
	logger.Info(`Handling update ID=%d, Message=%d`, options.Upd.UpdateId, options.Upd.Message.MessageId)
	logger.Debug(`> %s`, options.Upd.Message.Text)

	var hasProcessed = false

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
			}
		}
	}

	// todo: update index even if no text messages can be processed
	// TODO: always single send-message for message+message entities?

	if !hasProcessed && options.Upd.Message.Text != `` {
		msg, err := r.commandDefault.Run(options)
		if err != nil {
			logger.Error(`Failed to run command "%s": %s`, r.commandDefault.GetName(), err)
		} else {
			r.sendMessage(msg)
		}

		r.updateOffset(options) // do not stop on failed command
	} else if !hasProcessed {
		logger.Debug(`No new messages...`)
	} else {
		r.updateOffset(options)
	}
}

func (r *TelegramBotsApiStruct) updateOffset(options RunOptionsStruct) {
	if options.Upd.UpdateId >= r.routingUpdate.Offset {
		logger.Debug("Was offset %d, will be: %d", r.routingUpdate.Offset, options.Upd.UpdateId+1)
		r.routingUpdate.Offset = options.Upd.UpdateId + 1
	}
}

func (r *TelegramBotsApiStruct) sendMessage(sendMessage SendMessageStruct) {
	if logger.DebugLevel() {
		logger.Debug(`Message to send: %s`, encodeToJson(sendMessage))
	}

	if _, ok := r.RequestManager.SendPostJsonRequest(r.routingSend.Uri(), sendMessage); !ok {
		logger.Error("Failed to send message: %s", encodeToJson(sendMessage))
	}
}

func (r *TelegramBotsApiStruct) processMessageEntity(runOptions RunOptionsStruct) (sendMessage SendMessageStruct, found bool) {
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

/////////////////

func NewTelegramBotsApi(authKey string, sleep int) *TelegramBotsApiStruct {
	var requestManager RequestManagerStruct

	api := TelegramBotsApiStruct{
		BaseUri:        apiBaseUri,
		RequestManager: &requestManager,
		AuthKey:        authKey,
		Sleep:          sleep,
	}

	api.routingMe.Init(&api)
	api.routingUpdate.Init(&api)
	api.routingSend.Init(&api)

	api.commands = append(
		api.commands,
		NewStartBotCommand(),
		NewAddTaskBotCommand(),
		NewListTaskBotCommand(),
	)

	return &api
}
