package main

import "encoding/json"

//const U+1F320
const emojiGlowingStar = "\U0001F31F"
const emojiColdSweat = "\U0001F613"

type BotCommander interface {
	Run(options RunOptionsStruct) (sendMessageStruct, error)
	GetName() string
}

type RunOptionsStruct struct {
	Upd UpdateTelegramModel
	Ent MessageEntityTelegramModel
}

type BotCommand struct {
	//Name string
}

//func (c BotCommand) Run(options RunOptionsStruct) (sendMessageStruct, error) {
//	// TODO: implement me
//}

//
//func (c BotCommand) GetName() string {
//	return c.Name
//}

type StartBotCommandStruct struct {
	BotCommand
	BotCommander
}

func (c StartBotCommandStruct) Run(options RunOptionsStruct) (sendMessageStruct, error) {
	logger.Debug(`Running %s command`, c.GetName())

	if logger.DebugLevel() {
		if data, err := json.Marshal(options); err == nil {
			logger.Debug(`Input: %s`, data)
		} else {
			logger.Info(`Failed encoding to JSON options: %s`, err)
		}
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		`Welcome to RPG Tasks Bot chat. Gain lots of XP and earn achievements and reach the goals! `+emojiGlowingStar, options.Upd.Message.MessageId), nil
}

func (c StartBotCommandStruct) GetName() string {
	return `/start`
}

type DefaultBotCommandStruct struct {
	BotCommand
	BotCommander
}

func (c DefaultBotCommandStruct) GetName() string {
	return `/default`
}

func (c DefaultBotCommandStruct) Run(options RunOptionsStruct) (sendMessageStruct, error) {
	//if len(upd.Message.Text) > ent.Offset+ent.Length {
	//	text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
	//} else {
	//text = `Sorry, cannot process your command`
	//}

	logger.Debug(`Running %s command`, c.GetName())

	return NewSendMessage(options.Upd.Message.Chat.Id, `Sorry, could not recognize this action`+emojiColdSweat, options.Upd.Message.MessageId), nil
}
