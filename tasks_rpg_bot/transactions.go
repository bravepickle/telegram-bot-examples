package main

import (
	"fmt"
	"strconv"
)

type Transactional interface {
	GetName() string
	//// GetData returns transactional data
	GetData() map[string]interface{}
	SetDataValue(name string, value interface{})
	GetDataValue(name string, defaultValue interface{}) interface{}
	GetFlashValue(name string, defaultValue interface{}) interface{} // get value and delete it from data array
	SetFlashValue(name string, value interface{})                    // get value and delete it from data array
	DelDataValue(name string) bool                                   // return true if found and deleted
	Init()
	Current() TransactionalStep
	//GetSteps() []string
	Next() TransactionalStep
	Prev() TransactionalStep
	Reset()
	Run(options RunOptionsStruct) (sendMessageStruct, bool)
	RunNextStep(options RunOptionsStruct) (sendMessageStruct, bool)
	Restart(options RunOptionsStruct) (sendMessageStruct, bool)
	Complete(options RunOptionsStruct) (sendMessageStruct, bool) // call when everything is set and we want to finish transaction
	//Commit() bool
	// TODO: for each user -> chat store values
}

type TransactionalStep interface {
	GetName() string
	Run(t Transactional, options RunOptionsStruct) (sendMessageStruct, bool)
	Revert(t Transactional, options RunOptionsStruct) // revert run step
	//Set(name string, value interface{}) // set value for the step
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
		logger.Error(`Current step not defined: %d`, t.currentStepIndex)

		return nil, false
	}
	//return t.Current().Run(t, options)
	//t.RunStep()

	//return nil, true
	//return nil, false
}

func (t *TransactionStruct) RunNextStep(options RunOptionsStruct) (sendMessageStruct, bool) {
	nextStep := t.Next() // next step to do...

	if nextStep == nil {
		// TODO: save task data! transaction.Complete() or something similar : save + delete completed transaction

		return t.Complete(options)

		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`All data is filled: `+string(encodeToJson(t.GetData())), 0), true
	}

	options.Upd.Message.Text = `` // hack to start processing new step

	return nextStep.Run(t, options)
}

func (t *TransactionStruct) Restart(options RunOptionsStruct) (sendMessageStruct, bool) {
	t.Reset()
	t.Init()

	step := t.Current() // next step to do...

	if step == nil {
		// TODO: save task data! transaction.Complete() or something similar : save + delete completed transaction

		return t.Complete(options)

		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`All data is filled: `+string(encodeToJson(t.GetData())), 0), true
	}

	options.Upd.Message.Text = `` // hack to start processing new step

	return step.Run(t, options)
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

func (t *TransactionStruct) Complete(options RunOptionsStruct) (sendMessageStruct, bool) {

	return NewSendMessage(options.Upd.Message.Chat.Id, "TBD: transaction is completed `"+string(encodeToJson(t.GetData()))+"`", 0), true
	//if value, ok := t.data[name]; ok {
	//	return value
	//}
	//
	//return defaultValue
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
	t.steps = append(t.steps, &SummaryStep{}) // TODO: add mapping for fields or use toString in steps to convert
	t.steps = append(t.steps, &ConfirmStep{})
	t.steps = append(t.steps, &TaskDefaultStep{})
}

//func (t *AddTaskTransactionStruct) Run() bool {
//	t.RunStep()
//
//	return true
//}
// =========== BasicStep
//type BasicStep struct {
//}
//
//func (s *BasicStep) Set(name string, value interface{}) {
//	// do nothing. override in dependencies
//}

// =========== ExperienceStep
//"CREATE TABLE IF NOT EXISTS task (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, title TEXT, description TEXT, status TEXT, exp INTEGER, date_expiration TEXT DEFAULT '', date_created TEXT DEFAULT CURRENT_TIMESTAMP, date_updated TEXT DEFAULT CURRENT_TIMESTAMP)",
type ExperienceStep struct {
	//BasicStep
}

func (t ExperienceStep) GetName() string {
	return `experience`
}

func (t ExperienceStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	if options.Upd.Message.Entities == nil && options.Upd.Message.Text != `` {
		//task := tr.GetDataValue(`task`, &TaskDbEntity{})
		task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)

		var err error

		task.Exp, err = strconv.Atoi(options.Upd.Message.Text)

		if err == nil {
			tr.SetDataValue(`task`, task)

			return tr.RunNextStep(options)
		}

		logger.Info(`Failed to validate exp value: %s`, err)
		//tr.SetDataValue(`exp`, options.Upd.Message.Text)

		//text := tr.GetDataValue(`exp`, ``).(string)

		//tr.Next() // next step to do...

		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`Your experience will be: ` + text, options.Upd.Message.MessageId), true
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		`Please, enter amount of experience gained for the completion task`, 0), true
}

func (t ExperienceStep) Revert(tr Transactional, options RunOptionsStruct) {
	//tr.SetDataValue(`title`, nil) // TODO: revert properly - check state and decide what to do
	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.Exp = 0
	tr.SetDataValue(`task`, task)
}

// =========== TitleStep
type TitleStep struct {
	//BasicStep
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
		//tr.SetDataValue(`title`, options.Upd.Message.Text)

		task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
		task.Title = options.Upd.Message.Text
		tr.SetDataValue(`task`, task)

		return tr.RunNextStep(options)

		//nextStep := tr.Next() // next step to do...
		//
		//if nextStep == nil {
		//	// TODO: properly process here. What should we do at the end?
		//	//return nil, true
		//
		//	return NewSendMessage(options.Upd.Message.Chat.Id,
		//		`All data is filled: `+string(encodeToJson(tr.GetData())), 0), true
		//}
		//
		//options.Upd.Message.Text = `` // hack to start processing new step
		//
		//return nextStep.Run(tr, options)

		//text := tr.GetDataValue(`title`, ``).(string)
		//
		//return NewSendMessage(options.Upd.Message.Chat.Id,
		//	`Your title will be: ` + text, options.Upd.Message.MessageId), true
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		`Please, enter title for the task`, 0), true

	//return true
}

func (t TitleStep) Revert(tr Transactional, options RunOptionsStruct) {
	//tr.SetDataValue(`title`, nil) // TODO: revert properly - check state and decide what to do

	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.Title = ``
	tr.SetDataValue(`task`, task)
}

// =========== TaskDefaultStep
type TaskDefaultStep struct {
	//BasicStep
}

func (t TaskDefaultStep) GetName() string {
	return `task-default`
}

func (t TaskDefaultStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	// TODO: set task db entity instead - more handy!
	//tr.SetDataValue(`user_id`, options.Upd.Message.From.Id)
	//tr.SetDataValue(`status`, statusPending)

	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.UserId = int(options.Upd.Message.From.Id)
	task.Status = statusPending
	tr.SetDataValue(`task`, task)

	return tr.RunNextStep(options)

	//Id             int
	//UserId         int
	//Title          string
	//Status         string
	//Exp            int
	//Description    string
	//DateExpiration string

	return tr.RunNextStep(options)

	//nextStep := tr.Next() // next step to do...
	//
	//if nextStep == nil {
	//	// TODO: save task data! transaction.Complete() or something similar : save + delete completed transaction
	//
	//	return tr.Complete(options)
	//
	//	//return NewSendMessage(options.Upd.Message.Chat.Id,
	//	//	`All data is filled: `+string(encodeToJson(tr.GetData())), 0), true
	//}
	//
	//options.Upd.Message.Text = `` // hack to start processing new step
	//
	//return nextStep.Run(tr, options)

}

func (t TaskDefaultStep) Revert(tr Transactional, options RunOptionsStruct) {
	//tr.SetDataValue(`user_id`, 0)
	//tr.SetDataValue(`status`, ``)

	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.UserId = 0
	task.Status = ``
	tr.SetDataValue(`task`, task)
}

// =========== SummaryStep
type SummaryStep struct {
	//BasicStep

	Shown bool
}

func (t SummaryStep) GetName() string {
	return `summary`
}

//func (t *SummaryStep) SetShown(value bool) {
//	t.Shown = value
//}

func (t *SummaryStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	//shown := tr.GetDataValue(`is_summary_shown`, false).(bool)
	//shown := t.Shown
	logger.Debug(`>>>>>>>>>>>>>>>>>>>>>>>>> Task "%s" shown = %t`, t.GetName(), t.Shown)

	if !t.Shown {
		t.Shown = true
		//tr.SetDataValue(`is_summary_shown`, true)
		text := "*Summary:* \n"
		data := tr.GetData()

		logger.Debug(`+++ >>>>>>>>>>>>>>>>>>>>>>>>> Task "%s" shown = %t`, t.GetName(), t.Shown)

		// TODO: use params mapping and values type check

		for name, value := range data {

			switch valType := value.(type) {
			case DbEntityInterface:
				text += fmt.Sprintf("  %s: `%s`\n", name, encodeToJson(value))
			default:
				text += fmt.Sprintf("  %s: (%s) `%v`\n", name, valType, value)

			}

			//text += fmt.Sprintf("  %s: `(%T) [%v] %v`\n", name, value, value.(DbEntityInterface), value)
		}

		tr.SetDataValue(`message_text_prepend`, text) // a hack to avoid sending message right away. Will receive in confirm msg

		//return NewSendMessage(options.Upd.Message.Chat.Id, text, 0), true
	}

	//Id             int
	//UserId         int
	//Title          string
	//Status         string
	//Exp            int
	//Description    string
	//DateExpiration string

	return tr.RunNextStep(options)

	//nextStep := tr.Next() // next step to do...
	//
	//if nextStep == nil {
	//	// TODO: save task data! transaction.Complete() or something similar : save + delete completed transaction
	//
	//	return tr.Complete(options)
	//
	//	//return NewSendMessage(options.Upd.Message.Chat.Id,
	//	//	`All data is filled: `+string(encodeToJson(tr.GetData())), 0), true
	//}
	//
	//options.Upd.Message.Text = `` // hack to start processing new step
	//
	//return nextStep.Run(tr, options)

}

func (t *SummaryStep) Revert(tr Transactional, options RunOptionsStruct) {
	t.Shown = false
	//tr.SetDataValue(`message_text_prepend`, ``)
	tr.DelDataValue(`message_text_prepend`)
	//tr.SetDataValue(`is_summary_shown`, false)
}

// =========== ConfirmStep
type ConfirmStep struct {
	TransactionalStep

	Shown bool
}

func (t ConfirmStep) GetName() string {
	return `summary`
}

func (t *ConfirmStep) SetShown(value bool) {
	t.Shown = value
}

func (t *ConfirmStep) Run(tr Transactional, options RunOptionsStruct) (sendMessageStruct, bool) {
	var yes, no = `y`, `n`
	if !t.Shown {
		t.Shown = true
		//tr.Current().Se = true
		//t.SetShown(true)

		// TODO: on typing NO go to step 1

		prependText := tr.GetFlashValue(`message_text_prepend`, ``).(string)

		var text string

		if prependText != `` {
			//tr.DelDataValue(`message_text_prepend`)
			//tr.SetDataValue(`message_text_prepend`, ``) // TODO: add DelDataValue instead

			text = fmt.Sprintf("%s\n*Proceed?* %s/%s", prependText, yes, no)
		} else {
			text = fmt.Sprintf("Proceed? %s/%s", yes, no)
		}

		return NewSendMessage(options.Upd.Message.Chat.Id, text, 0), true
	}

	if options.Upd.Message.Text == no {
		return tr.Restart(options)
	}

	//Id             int
	//UserId         int
	//Title          string
	//Status         string
	//Exp            int
	//Description    string
	//DateExpiration string

	return tr.RunNextStep(options)

	//nextStep := tr.Next() // next step to do...
	//
	//if nextStep == nil {
	//	// TODO: save task data! transaction.Complete() or something similar : save + delete completed transaction
	//
	//	return tr.Complete(options)
	//
	//	//return NewSendMessage(options.Upd.Message.Chat.Id,
	//	//	`All data is filled: `+string(encodeToJson(tr.GetData())), 0), true
	//}
	//
	//options.Upd.Message.Text = `` // hack to start processing new step
	//
	//return nextStep.Run(tr, options)

}

func (t ConfirmStep) Revert(tr Transactional, options RunOptionsStruct) {
	t.Shown = false
}

// ===========

func NewAddTaskTransaction() *AddTaskTransactionStruct {
	var transaction AddTaskTransactionStruct
	transaction.Init()

	return &transaction
}

func NewStartBotCommand() (model StartBotCommandStruct) {
	//model.Init()

	return model
}
func NewAddTaskBotCommand() (model AddTaskBotCommandStruct) {
	model.transactions = make(map[uint32]map[uint32]Transactional) // TODO: fix it somehow, see https://stackoverflow.com/questions/40823315/go-x-does-not-implement-y-method-has-a-pointer-receiver
	//model.Init()

	return model
}
