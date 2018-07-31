package main

type Transactional interface {
	GetName() string
	//// GetData returns transactional data
	GetData() map[string]interface{}
	SetDataValue(name string, value interface{})
	GetDataValue(name string, defaultValue interface{}) interface{}
	Init()
	Current() TransactionalStep
	//GetSteps() []string
	Next() TransactionalStep
	Prev() TransactionalStep
	Reset()
	Run(options RunOptionsStruct) (sendMessageStruct, bool)
	//Complete() bool
	//Commit() bool
	// TODO: for each user -> chat store values
}

type TransactionalStep interface {
	GetName() string
	Run(t Transactional, options RunOptionsStruct) (sendMessageStruct, bool)
	Revert(t Transactional, options RunOptionsStruct) // revert run step
}

// =========== TransactionStruct
type TransactionStruct struct {
	currentStepIndex int
	//steps            map[int]TransactionalStep
	steps []TransactionalStep
	data  map[string]interface{}
}

func (t *TransactionStruct) Next() TransactionalStep {
	t.currentStepIndex += 1

	return t.Current()
}

func (t *TransactionStruct) Prev() TransactionalStep {
	t.currentStepIndex -= 1

	return t.Current()
}

func (t *TransactionStruct) Reset() {
	t.currentStepIndex = 0
	t.data = make(map[string]interface{})
	//t.steps = make(map[int]TransactionalStep)
	//t.steps = make(map[int]TransactionalStep)
}

func (t *TransactionStruct) Current() TransactionalStep {
	if len(t.steps) <= t.currentStepIndex {
		//logger.Error(`No steps initialized for transaction: %s[%d]`, t.GetName(), t.currentStepIndex)

		//logger.Error(`STEP: %v`, step)
		//logger.Error(`STEP JSON: %s`, encodeToJson(step))

		return nil
	} else {
		return t.steps[t.currentStepIndex]
	}
}

func (t *TransactionStruct) Init() {
}

func (t *TransactionStruct) Run(options RunOptionsStruct) (sendMessageStruct, bool) {
	currentStep := t.Current()
	if currentStep != nil {
		return currentStep.Run(t, options)
	} else {
		logger.Error(`Current step not defined: %s`, t.GetName())

		return nil, false
	}
	//return t.Current().Run(t, options)
	//t.RunStep()

	//return nil, true
	//return nil, false
}

//func (t *TransactionStruct) RunStep() {
//	//if step, ok := t.steps[t.currentStepIndex]; !ok {
//	//if step := t.steps[t.currentStepIndex]; !ok {
//	if step := t.steps[t.currentStepIndex]; step != nil {
//		logger.Error(`No steps initialized for transaction %s`, t.GetName())
//	} else {
//		step.Run(t)
//	}
//}

func (t *TransactionStruct) GetName() string {
	return `[undefined]` // override in children
}

func (t *TransactionStruct) GetData() map[string]interface{} {
	return t.data
}

func (t *TransactionStruct) SetDataValue(name string, value interface{}) {
	t.data[name] = value
}

func (t *TransactionStruct) GetDataValue(name string, defaultValue interface{}) interface{} {
	if value, ok := t.data[name]; ok {
		return value
	}

	return defaultValue
}

// =========== AddTaskTransactionStruct

type AddTaskTransactionStruct struct {
	TransactionStruct
}

func (t *AddTaskTransactionStruct) GetName() string {
	return `add-task`
}

func (t *AddTaskTransactionStruct) Init() {
	t.Reset()

	t.steps = append(t.steps, TitleStep{})
	t.steps = append(t.steps, ExperienceStep{})
}

//func (t *AddTaskTransactionStruct) Run() bool {
//	t.RunStep()
//
//	return true
//}

// =========== ExperienceStep
//"CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, title TEXT, description TEXT, status TEXT, exp INTEGER, date_expiration TEXT DEFAULT '', date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)",
type ExperienceStep struct {
	TransactionalStep
}

func (t ExperienceStep) GetName() string {
	return `experience`
}

func (t ExperienceStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	if options.Upd.Message.Entities == nil && options.Upd.Message.Text != `` {
		tr.SetDataValue(`exp`, options.Upd.Message.Text)

		nextStep := tr.Next() // next step to do...

		if nextStep == nil {
			// TODO: properly process here. What should we do at the end?
			//return nil, true

			return NewSendMessage(options.Upd.Message.Chat.Id,
				`All data is filled: `+string(encodeToJson(tr.GetData())), options.Upd.Message.MessageId), true
		}

		options.Upd.Message.Text = `` // hack to start processing new step

		return nextStep.Run(tr, options)

		//text := tr.GetDataValue(`exp`, ``).(string)

		//tr.Next() // next step to do...

		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`Your experience will be: ` + text, options.Upd.Message.MessageId), true
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		`Please, enter amount of experience gained for the completion task`, options.Upd.Message.MessageId), true
}

func (t ExperienceStep) Revert(tr Transactional, options RunOptionsStruct) {
	tr.SetDataValue(`title`, nil) // TODO: revert properly - check state and decide what to do
}

// =========== TitleStep
type TitleStep struct {
	TransactionalStep
}

func (t TitleStep) GetName() string {
	return `title`
}

func (t TitleStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	// TODO: generate message here to input
	//tr.SetData(`title`, `abcddsd // `)

	//tr.SetDataValue(`title`, nil)

	//if options.Upd.Message.Chat.Id == 0 {
	//	logger.Error(`Failed to find chat ID for update message: %s`, encodeToJson(options.Upd))
	//
	//	return nil, false
	//}

	//logger.Info(`>>>>> Reading title or prompt? %s`, encodeToJson(options))

	if options.Upd.Message.Entities == nil && options.Upd.Message.Text != `` {
		tr.SetDataValue(`title`, options.Upd.Message.Text)

		nextStep := tr.Next() // next step to do...

		if nextStep == nil {
			// TODO: properly process here. What should we do at the end?
			//return nil, true

			return NewSendMessage(options.Upd.Message.Chat.Id,
				`All data is filled: `+string(encodeToJson(tr.GetData())), options.Upd.Message.MessageId), true
		}

		options.Upd.Message.Text = `` // hack to start processing new step

		return nextStep.Run(tr, options)

		//text := tr.GetDataValue(`title`, ``).(string)
		//
		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`Your title will be: ` + text, options.Upd.Message.MessageId), true
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		`Please, enter title for the task`, options.Upd.Message.MessageId), true

	//return true
}

func (t TitleStep) Revert(tr Transactional, options RunOptionsStruct) {
	tr.SetDataValue(`title`, nil) // TODO: revert properly - check state and decide what to do
}

// ===========

func NewAddTaskTransaction() *AddTaskTransactionStruct {
	var transaction AddTaskTransactionStruct
	transaction.Init()

	return &transaction
}

func NewStartBotCommand() (model StartBotCommandStruct) {
	model.Init()

	return model
}
func NewAddTaskBotCommand() (model AddTaskBotCommandStruct) {
	model.transactions = make(map[uint32]map[uint32]Transactional) // TODO: fix it somehow, see https://stackoverflow.com/questions/40823315/go-x-does-not-implement-y-method-has-a-pointer-receiver
	model.Init()

	return model
}
