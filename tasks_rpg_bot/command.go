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
}

type RunOptionsStruct struct {
	Upd UpdateTelegramModel
	Ent MessageEntityTelegramModel
}

type BotCommand struct {
}

func (c BotCommand) IsRunning(options RunOptionsStruct) bool {
	return false
}

// ================== StartBotCommandStruct

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

	return NewSendMessage(options.Upd.Message.Chat.Id, usrMsg.T(`response.welcome`), options.Upd.Message.MessageId), nil
}

func (c StartBotCommandStruct) GetName() string {
	return `/start`
}

// ================== DefaultBotCommandStruct

type DefaultBotCommandStruct struct {
	BotCommand
}

func (c DefaultBotCommandStruct) GetName() string {
	return `/default`
}

func (c DefaultBotCommandStruct) Run(options RunOptionsStruct) (SendMessageStruct, error) {
	logger.Debug(`Running %s command`, c.GetName())

	return NewSendMessage(options.Upd.Message.Chat.Id, usrMsg.T(`response.unrecognized`), options.Upd.Message.MessageId), nil
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

	trans := c.initTransaction(options)
	if sendMessage, ok := trans.Run(options); ok {
		if logger.DebugLevel() {
			logger.Debug(`Transaction step data result: %s`, encodeToJson(sendMessage))
		}

		return sendMessage, nil
	} else {
		return nil, errors.New(`Failed to run command.`)
	}
}

func (c AddTaskBotCommandStruct) IsRunning(options RunOptionsStruct) bool {
	chatId := options.Upd.Message.Chat.Id
	userId := options.Upd.Message.From.Id

	_, ok := c.transactions[chatId][userId]

	return ok
}

// ================== ListTaskBotCommandStruct

type ListTaskBotCommandStruct struct {
	BotCommand
}

func (c ListTaskBotCommandStruct) Run(options RunOptionsStruct) (SendMessageStruct, error) {
	logger.Debug(`Running %s command`, c.GetName())

	msg := usrMsg.T(`response.task.list_header`) + "\n"
	tasks := dbManager.findTasksByUser(int(options.Upd.Message.Chat.Id))

	if len(tasks) == 0 {
		msg += usrMsg.T(`response.task.list_item`)
	} else {
		for k, task := range tasks {
			//%d. _%s_: expires at "%s", gain "%d" XP
			msg += usrMsg.T(`response.task.list_item`, k+1, task.Title, task.DateExpiration, task.Exp) + "\n"
		}
	}

	return NewSendMessage(options.Upd.Message.Chat.Id, msg, 0), nil
}

func (c ListTaskBotCommandStruct) GetName() string {
	return `/list`
}

// ==================
