package main

import (
	"encoding/json"
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
func (r TelegramBotsApiRequestModel) Uri() string {
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

type TelegramBotsApiStruct struct {
	BaseUri        string                // base API URI
	RequestManager *RequestManagerStruct // handling requests
	AuthKey        string                // API auth key
	BotInfo        BotProfileTelegramModel

	routingMe MeRequestModel

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

/////////////////

func NewTelegramBotsApi(authKey string) *TelegramBotsApiStruct {
	var requestManager RequestManagerStruct

	api := TelegramBotsApiStruct{
		BaseUri:        apiBaseUri,
		RequestManager: &requestManager,
		AuthKey:        authKey,
	}

	api.routingMe.init(&api)

	return &api
}
