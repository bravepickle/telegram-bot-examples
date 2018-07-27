package main

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
	logger.Info(`Running %s command`, c.GetName())

	return NewSendMessage(options.Upd.Message.Chat.Id, `Results of running command `+c.GetName()), nil
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

	logger.Info(`Running %s command`, c.GetName())

	return NewSendMessage(options.Upd.Message.Chat.Id, `Results of running command `+c.GetName()), nil
}
