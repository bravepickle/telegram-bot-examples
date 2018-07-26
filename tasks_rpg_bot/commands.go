package main

type BotCommander interface {
	Run(options RunOptionsStruct)
	GetName() string
}

type RunOptionsStruct struct {
	Upd UpdateTelegramModel
	Ent MessageEntityTelegramModel
}

type BotCommand struct {
	Name string
}

func (c BotCommand) Run(options RunOptionsStruct) {
	// TODO: implement me
}

func (c BotCommand) GetName() string {
	return c.Name
}

type StartBotCommandStruct struct {
	BotCommand
	//BotCommander
}

func (c StartBotCommandStruct) Run(options RunOptionsStruct) {
	// TODO: implement me
}

type DefaultBotCommandStruct struct {
	BotCommand
	//BotCommander
}

func (c DefaultBotCommandStruct) Run(options RunOptionsStruct) {
	//if len(upd.Message.Text) > ent.Offset+ent.Length {
	//	text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
	//} else {
	text = `Sorry, cannot process your command`
	//}
}