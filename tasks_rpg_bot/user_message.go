package main

import "fmt"

const defaultUnrecognizedMsg = "*[[!%s!]]*"

type UserMessage struct {
	messages map[string]string
}

// T search message by alias and return it with arguments formatted
func (m *UserMessage) T(name string, args ...interface{}) string {
	if msg, ok := m.messages[name]; ok {
		return fmt.Sprintf(msg, args...)
	}

	logger.Error(`Failed to find user message for: %s`, name)

	return fmt.Sprintf(defaultUnrecognizedMsg, name)
}

func NewUserMessage() *UserMessage {
	msg := UserMessage{}
	msg.messages = map[string]string{
		`answer.yes`: `yes`,
		`answer.no`:  `no`,

		`remind.task.header`: `Tasks TODO:`,
		`remind.task.line`:   `%d. _%s_: expires at "%s", gain "%d" XP`,

		`response.unrecognized`:   `Sorry, could not recognize this action` + emojiColdSweat,
		`response.welcome`:        `Welcome to RPG Tasks Bot chat. Gain lots of XP and earn achievements and reach the goals! ` + emojiGlowingStar,
		`response.save.fail`:      `Failed to save data. Please, try again lager`,
		`response.save.success`:   "TBD: transaction is completed `%s`",
		`response.date.fail`:      `Failed to read date. Please, try again`,
		`response.summary.header`: `Summary`,

		`request.proceed`:         `*Proceed?* %s/%s`,
		`request.task.title`:      `Please, enter title for the task`,
		`request.task.exp`:        `Please, enter amount of experience gained for the completion task`,
		`request.task.expiration`: `Please, enter expiration date for the task (formats: "%s", "%s") or write "%s"`,
	}

	return &msg
}
