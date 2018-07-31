package main

type Transactional interface {
	GetName() string
	//// GetData returns transactional data
	GetData() map[string]interface{}
	SetData(name string, value interface{})
	Init()
	Current() TransactionalStep
	//GetSteps() []string
	Next()
	Prev()
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

func (t *TransactionStruct) Next() {
	t.currentStepIndex += 1
}

func (t *TransactionStruct) Prev() {
	t.currentStepIndex -= 1
}

func (t *TransactionStruct) Reset() {
	t.currentStepIndex = 0
	t.data = make(map[string]interface{})
	//t.steps = make(map[int]TransactionalStep)
	//t.steps = make(map[int]TransactionalStep)
}

func (t *TransactionStruct) Current() TransactionalStep {
	if step := t.steps[t.currentStepIndex]; step == nil {
		logger.Error(`No steps initialized for transaction: %s[%d]`, t.GetName(), t.currentStepIndex)

		//logger.Error(`STEP: %v`, step)
		//logger.Error(`STEP JSON: %s`, encodeToJson(step))

		return nil
	} else {
		return step
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

func (t *TransactionStruct) SetData(name string, value interface{}) {
	t.data[name] = value
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

	t.steps = append(t.steps, SetTitleStep{})
}

//func (t *AddTaskTransactionStruct) Run() bool {
//	t.RunStep()
//
//	return true
//}

// =========== SetTitleStep
type SetTitleStep struct {
	TransactionalStep
}

func (t SetTitleStep) GetName() string {
	return `set-title`
}

func (t SetTitleStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	// TODO: generate message here to input
	//tr.SetData(`title`, `abcddsd // `)

	if options.Upd.Message.Chat.Id == 0 {
		logger.Error(`Failed to find chat ID for update message: %s`, encodeToJson(options.Upd))

		return nil, false
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		`Please, enter title for the task`, options.Upd.Message.MessageId), true

	//return true
}
func (t SetTitleStep) Revert(tr Transactional, options RunOptionsStruct) {
	tr.SetData(`title`, nil) // TODO: revert properly - check state and decide what to do
}

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
	model.Init()

	return model
}
