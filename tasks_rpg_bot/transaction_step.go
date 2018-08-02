package main

import (
	"fmt"
	"strconv"
	"time"
)

type TransactionalStep interface {
	GetName() string
	Run(t Transactional, options RunOptionsStruct) (SendMessageStruct, bool)
	Revert(t Transactional, options RunOptionsStruct) // revert run step
	//Set(name string, value interface{}) // set value for the step
}

// =========== TitleStep
type TitleStep struct{}

func (t TitleStep) GetName() string {
	return `title`
}

func (t TitleStep) Run(tr Transactional, options RunOptionsStruct) (SendMessageStruct, bool) {
	if options.Upd.Message.Entities == nil && options.Upd.Message.Text != `` {
		task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
		task.Title = options.Upd.Message.Text
		tr.SetDataValue(`task`, task)

		return tr.RunNextStep(options)
	}

	if options.Ent.Length > 0 && len(options.Upd.Message.Text) > options.Ent.Length {
		task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
		task.Title = options.Upd.Message.Text[options.Ent.Length+1:]
		tr.SetDataValue(`task`, task)

		return tr.RunNextStep(options)
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		usrMsg.T(`request.task.title`), 0), true
}

func (t TitleStep) Revert(tr Transactional, options RunOptionsStruct) {
	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.Title = ``
	tr.SetDataValue(`task`, task)
}

// =========== ExperienceStep
type ExperienceStep struct{}

func (t ExperienceStep) GetName() string {
	return `experience`
}

func (t ExperienceStep) Run(tr Transactional, options RunOptionsStruct) (SendMessageStruct, bool) {
	if options.Upd.Message.Entities == nil && options.Upd.Message.Text != `` {
		task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)

		var err error
		task.Exp, err = strconv.Atoi(options.Upd.Message.Text)

		if err == nil {
			tr.SetDataValue(`task`, task)

			return tr.RunNextStep(options)
		}

		logger.Info(`Failed to validate exp value: %s`, err)
	}

	return NewSendMessage(options.Upd.Message.Chat.Id,
		usrMsg.T(`request.task.exp`), 0), true
}

func (t ExperienceStep) Revert(tr Transactional, options RunOptionsStruct) {
	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.Exp = 0
	tr.SetDataValue(`task`, task)
}

// =========== DateExpirationStep
type DateExpirationStep struct{}

func (t DateExpirationStep) GetName() string {
	return `date-expiration`
}

func (t DateExpirationStep) Run(tr Transactional, options RunOptionsStruct) (SendMessageStruct, bool) {
	answerNo := usrMsg.T(`answer.no`)

	if options.Upd.Message.Entities == nil && options.Upd.Message.Text != `` {
		if answerNo != options.Upd.Message.Text {
			dt, err := time.Parse("02/01", options.Upd.Message.Text)

			// TODO: add validation date in future
			if err != nil {
				dt, err = time.Parse("2006-01-02", options.Upd.Message.Text)
				if err != nil {
					logger.Error(`Failed to parse date "%s": %s`, options.Upd.Message.Text, err)

					return NewSendMessage(options.Upd.Message.Chat.Id, usrMsg.T(`response.date.fail`), 0), true
				}
			} else {
				dt = dt.AddDate(time.Now().Year(), 0, 0) // append current year

				if dt.Before(time.Now()) {
					dt = dt.AddDate(1, 0, 0) // append next year if date was already exceeded
				}
			}

			task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
			task.DateExpiration = NewDbTime(dt)
			tr.SetDataValue(`task`, task)
		}

		return tr.RunNextStep(options)
	}

	text := usrMsg.T(`request.task.expiration`, dateFormat, dateFormatShort, answerNo)

	return NewSendMessage(options.Upd.Message.Chat.Id, text, 0), true
}

func (t DateExpirationStep) Revert(tr Transactional, options RunOptionsStruct) {
	//task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	//task.DateExpiration = ``
	//tr.SetDataValue(`task`, task)
}

// =========== TaskDefaultStep
type TaskDefaultStep struct{}

func (t TaskDefaultStep) GetName() string {
	return `task-default`
}

func (t TaskDefaultStep) Run(tr Transactional, options RunOptionsStruct) (SendMessageStruct, bool) {
	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.UserId = int(options.Upd.Message.From.Id)
	task.Status = statusPending
	task.DateUpdated = NewDbTime(time.Now())
	task.DateCreated = NewDbTime(time.Now())
	tr.SetDataValue(`task`, task)

	return tr.RunNextStep(options)
}

func (t TaskDefaultStep) Revert(tr Transactional, options RunOptionsStruct) {
	task := tr.GetDataValue(`task`, &TaskDbEntity{}).(*TaskDbEntity)
	task.UserId = 0
	task.Status = ``
	tr.SetDataValue(`task`, task)
}

// =========== SummaryStep
type SummaryStep struct {
	Shown bool
}

func (t SummaryStep) GetName() string {
	return `summary`
}

func (t *SummaryStep) Run(tr Transactional, options RunOptionsStruct) (SendMessageStruct, bool) {
	if !t.Shown {
		t.Shown = true
		text := "*" + usrMsg.T(`response.summary.header`) + ":* \n"
		data := tr.GetData()

		for name, value := range data {
			switch valType := value.(type) {
			case DbEntityInterface:
				text += fmt.Sprintf("  %s: `%s`\n", name, encodeToJson(value))
			default:
				text += fmt.Sprintf("  %s: (%s) `%v`\n", name, valType, value)

			}
		}

		tr.SetDataValue(`message_text_prepend`, text) // a hack to avoid sending message right away. Will receive in confirm msg
	}

	return tr.RunNextStep(options)
}

func (t *SummaryStep) Revert(tr Transactional, options RunOptionsStruct) {
	t.Shown = false
	tr.DelDataValue(`message_text_prepend`)
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

func (t *ConfirmStep) Run(tr Transactional, options RunOptionsStruct) (SendMessageStruct, bool) {
	var yes, no = usrMsg.T(`answer.yes`), usrMsg.T(`answer.no`)
	if !t.Shown {
		t.Shown = true
		prependText := tr.GetFlashValue(`message_text_prepend`, ``).(string)
		var text string

		if prependText != `` {
			text = prependText + "\n" + usrMsg.T(`request.proceed`, yes, no)
		} else {
			text = usrMsg.T(`request.proceed`, yes, no)
		}

		return NewSendMessage(options.Upd.Message.Chat.Id, text, 0), true
	}

	if options.Upd.Message.Text == no {
		return tr.Restart(options)
	}

	return tr.RunNextStep(options)
}

func (t ConfirmStep) Revert(tr Transactional, options RunOptionsStruct) {
	t.Shown = false
}
