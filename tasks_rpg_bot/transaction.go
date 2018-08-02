package main

const dateFormat = `YYYY-mm-dd`
const dateFormatShort = `dd/mm`
const answerNo = "no"

type Transactional interface {
	GetName() string
	GetData() map[string]interface{}
	SetDataValue(name string, value interface{})
	GetDataValue(name string, defaultValue interface{}) interface{}
	GetFlashValue(name string, defaultValue interface{}) interface{} // get value and delete it from data array
	SetFlashValue(name string, value interface{})                    // get value and delete it from data array
	DelDataValue(name string) bool                                   // return true if found and deleted
	Init()
	Current() TransactionalStep
	Next() TransactionalStep
	Prev() TransactionalStep
	Reset()
	Run(options RunOptionsStruct) (SendMessageStruct, bool)
	RunNextStep(options RunOptionsStruct) (SendMessageStruct, bool)
	Restart(options RunOptionsStruct) (SendMessageStruct, bool)
	Complete(options RunOptionsStruct) (SendMessageStruct, bool) // call when everything is set and we want to finish transaction
	// TODO: for each user -> chat store values
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

func (t *TransactionStruct) Run(options RunOptionsStruct) (SendMessageStruct, bool) {
	currentStep := t.Current()
	if currentStep != nil {
		return currentStep.Run(t, options)
	} else {
		logger.Error(`Current step not defined: %d`, t.currentStepIndex)

		return nil, false
	}
}

func (t *TransactionStruct) RunNextStep(options RunOptionsStruct) (SendMessageStruct, bool) {
	nextStep := t.Next()
	if nextStep == nil {
		return t.Complete(options)
	}

	options.Upd.Message.Text = `` // hack to start processing new step

	return nextStep.Run(t, options)
}

func (t *TransactionStruct) Restart(options RunOptionsStruct) (SendMessageStruct, bool) {
	t.Reset()
	t.Init()
	step := t.Current()
	if step == nil {
		return t.Complete(options)
	}

	options.Upd.Message.Text = `` // hack to start processing new step

	return step.Run(t, options)
}

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

func (t *TransactionStruct) DelDataValue(name string) bool {
	if _, ok := t.data[name]; ok {
		delete(t.data, name)

		return true
	}

	return false
}

func (t *TransactionStruct) SetFlashValue(name string, value interface{}) {
	t.SetDataValue(name, value)
}

func (t *TransactionStruct) GetFlashValue(name string, defaultValue interface{}) interface{} {
	value := t.GetDataValue(name, nil)

	if value == nil {
		return defaultValue
	}

	t.DelDataValue(name) // defer?

	return value
}

func (t *TransactionStruct) Complete(options RunOptionsStruct) (SendMessageStruct, bool) {
	for _, value := range t.GetData() {
		switch value.(type) {
		case DbEntityInterface:
			entity := value.(DbEntityInterface)
			if !entity.Save() {
				t.Reset()
				logger.Error(`Failed to save data of "%T" to DB: %s`, entity, encodeToJson(entity))

				return NewSendMessage(options.Upd.Message.Chat.Id, `Failed to save data. Please, try again lager`, 0), true
			}

		default:
			// do nothing
		}
	}

	text := "TBD: transaction is completed `" + string(encodeToJson(t.GetData())) + "`"

	t.Reset()

	return NewSendMessage(options.Upd.Message.Chat.Id, text, 0), true
}

// =========== AddTaskTransactionStruct
type AddTaskTransactionStruct struct {
	TransactionStruct
}

func (t *AddTaskTransactionStruct) GetName() string {
	return `add-task`
}

func (t *AddTaskTransactionStruct) Init() {
	logger.Debug(`>>>>>>>>>>> Init task "%s"`, t.GetName())
	t.Reset()

	t.steps = append(t.steps, &TitleStep{})
	t.steps = append(t.steps, &ExperienceStep{})
	t.steps = append(t.steps, &DateExpirationStep{})
	t.steps = append(t.steps, &TaskDefaultStep{})
	t.steps = append(t.steps, &SummaryStep{}) // TODO: add mapping for fields or use toString in steps to convert
	t.steps = append(t.steps, &ConfirmStep{})
}

// ===========

func NewAddTaskTransaction() *AddTaskTransactionStruct {
	var transaction AddTaskTransactionStruct
	transaction.Init()

	return &transaction
}

func NewStartBotCommand() (model StartBotCommandStruct) {
	return model
}
func NewAddTaskBotCommand() (model AddTaskBotCommandStruct) {
	model.transactions = make(map[uint32]map[uint32]Transactional) // TODO: fix it somehow, see https://stackoverflow.com/questions/40823315/go-x-does-not-implement-y-method-has-a-pointer-receiver

	return model
}
