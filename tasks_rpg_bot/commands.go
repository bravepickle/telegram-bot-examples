package main

import (
	"encoding/json"
	"errors"
)

//const U+1F320
const emojiGlowingStar = "\U0001F31F"
const emojiColdSweat = "\U0001F613"

type BotCommander interface {
	Run(options RunOptionsStruct) (SendMessageStruct, error)
	GetName() string
	IsRunning(options RunOptionsStruct) bool // return true if command is transactional and in process of running, e.g. waiting for user input
	//Init()
}

type RunOptionsStruct struct {
	Upd UpdateTelegramModel
	Ent MessageEntityTelegramModel
}

type BotCommand struct {
	//Name string
}

func (c BotCommand) IsRunning(options RunOptionsStruct) bool {
	return false
}

//func (c BotCommand) Init() {
//}

type StartBotCommandStruct struct {
	BotCommand
}

func (c StartBotCommandStruct) Run(options RunOptionsStruct) (SendMessageStruct, error) {
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
	//BotCommander
}

func (c DefaultBotCommandStruct) GetName() string {
	return `/default`
}

func (c DefaultBotCommandStruct) Run(options RunOptionsStruct) (SendMessageStruct, error) {
	//if len(upd.Message.Text) > ent.Offset+ent.Length {
	//	text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
	//} else {
	//text = `Sorry, cannot process your command`
	//}

	logger.Debug(`Running %s command`, c.GetName())

	return NewSendMessage(options.Upd.Message.Chat.Id, `Sorry, could not recognize this action`+emojiColdSweat, options.Upd.Message.MessageId), nil
}

// ================== AddTaskBotCommandStruct

type AddTaskBotCommandStruct struct {
	BotCommand

	// list of transactions that are running and not finished: v[userId][chatId] = Transactional
	/** ALTERNATIVE FORMAT, but unknown how to convert this to JSON: type Key struct {
	    Path, Country string
	}
	hits := make(map[Key]int
		**/
	transactions map[uint32]map[uint32]Transactional
}

func (c AddTaskBotCommandStruct) GetName() string {
	return `/add`
}

//func (c AddTaskBotCommandStruct) Init() {
//	//logger.Debug(`!!! Init bot command %s`, c.GetName())
//	//c.transactions = make(map[uint32]map[uint32]Transactional)
//}

func (c AddTaskBotCommandStruct) initTransaction(options RunOptionsStruct) Transactional {

	//logger.Debug(`Transaction value: %T %v`, c.transactions, c.transactions)
	//c.transactions = make(map[uint32]Transactional)

	//c.transactions = make(map[uint32]map[uint32]Transactional)

	chatId := options.Upd.Message.Chat.Id
	userId := options.Upd.Message.From.Id

	if trans, ok := c.transactions[chatId][userId]; !ok { // check if set
		if _, ok := c.transactions[chatId]; !ok {
			c.transactions[chatId] = make(map[uint32]Transactional)
		}

		//c.transactions[chatId] := make()

		c.transactions[chatId][userId] = NewAddTaskTransaction()

		return c.transactions[chatId][userId]
	} else {
		return trans
	}
}

func (c AddTaskBotCommandStruct) Run(options RunOptionsStruct) (SendMessageStruct, error) {
	logger.Debug(`Running %s command`, c.GetName())

	if logger.DebugLevel() {
		if data, err := json.Marshal(options); err == nil {
			logger.Debug(`Input: %s`, data)
		} else {
			logger.Info(`Failed encoding to JSON options: %s`, err)
		}
	}

	//if c.initChannel() {
	//	go c.RunChannel(options) // start channel running once
	//}
	//
	//sendMessage := <- c.channel

	trans := c.initTransaction(options)

	if sendMessage, ok := trans.Run(options); ok {
		//trans.Next()

		if logger.DebugLevel() {
			logger.Debug(`Transaction step data result: %s`, encodeToJson(sendMessage))
		}

		return sendMessage, nil

		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`Adding new task. Please, enter the title`, options.Upd.Message.MessageId), nil

	} else {
		//return nil, errors.New(`Failed to run command. ` + string(encodeToJson(sendMessage)))
		return nil, errors.New(`Failed to run command.`)
	}

	// TODO: remove return error if never used in all commands

	//if c.isRunning {
	//
	//} else {
	//	// TODO: implement me
	//}

	// TODO: channels pool and check transactions

	//return sendMessage, nil
	//return NewSendMessage(options.Upd.Message.Chat.Id,
	//	`Adding new task. Please, enter the title`, options.Upd.Message.MessageId), nil
}

func (c AddTaskBotCommandStruct) IsRunning(options RunOptionsStruct) bool {
	chatId := options.Upd.Message.Chat.Id
	userId := options.Upd.Message.From.Id

	//if trans, ok := c.transactions[chatId][userId];

	_, ok := c.transactions[chatId][userId]

	return ok
}

// ==================
