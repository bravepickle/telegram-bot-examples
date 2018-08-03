package main

import "fmt"

const defaultUnrecognizedMsg = "*[[!%s!]]*"
const defaultLocale = `en_US`
const localeRu = `ru_RU`
const localeUs = `en_US`

//var supportedLocales = []string{localeUs, localeRu}

type UserMessage struct {
	Locale string
	// locale -> alias
	messages map[string]map[string]string
}

// T search message by alias and return it with arguments formatted
func (m *UserMessage) T(name string, args ...interface{}) string {
	if msg, ok := m.messages[m.Locale][name]; ok {
		return fmt.Sprintf(msg, args...)
	} else if m.Locale != defaultLocale {
		// fallback to default locale
		if msg, ok := m.messages[defaultLocale][name]; ok {
			return fmt.Sprintf(msg, args...)
		}
	}

	logger.Error(`Failed to find user message for: %s`, name)

	return fmt.Sprintf(defaultUnrecognizedMsg, name)
}

func NewUserMessage(locale string) *UserMessage {
	msg := UserMessage{Locale: locale}

	msg.messages = make(map[string]map[string]string)

	msg.messages[localeUs] = map[string]string{
		`answer.yes`: `yes`,
		`answer.no`:  `no`,

		`remind.task.header`: `Tasks TODO:`,
		`remind.task.line`:   `%d. _%s_: expires at "%s", gain "%d" XP`,

		`response.unrecognized`:     `Sorry, could not recognize this action ` + emojiColdSweat,
		`response.welcome`:          `Welcome to RPG Tasks bot chat. Gain lots of XP and earn achievements and reach the goals! ` + emojiGlowingStar,
		`response.save.fail`:        `Failed to save data. Please, try again lager`,
		`response.save.success`:     "TBD: transaction is completed `%s`",
		`response.date.fail`:        `Failed to read date. Please, try again`,
		`response.summary.header`:   `Summary`,
		`response.task.list_header`: `*Tasks TODO:*`,
		`response.task.list_empty`:  `No tasks found.`,
		`response.task.list_item`:   `%d. _%s_: expires at "%s", gain "%d" XP`,

		`request.proceed`:         `*Proceed?* %s/%s`,
		`request.task.title`:      `Please, enter title for the task`,
		`request.task.exp`:        `Please, enter amount of experience gained for the completion task`,
		`request.task.expiration`: `Please, enter expiration date for the task (formats: "%s", "%s") or write "%s"`,
	}

	msg.messages[localeRu] = map[string]string{
		`answer.yes`: `да`,
		`answer.no`:  `нет`,

		`remind.task.header`: `Сделать задачи:`,
		`remind.task.line`:   `%d. _%s_: сделаешь до "%s", получишь "%d" опыта`,

		`response.unrecognized`:     `Не могу распознать команду` + emojiColdSweat,
		`response.welcome`:          `Добро пожаловать в RPG Tasks бот чат. Заработай ачивки, опыт и достигни своих целей! ` + emojiGlowingStar,
		`response.save.fail`:        `Не получилось сохранить данные. Попробуйте позже`,
		`response.save.success`:     "TBD: операция прошла успешно `%s`",
		`response.date.fail`:        `Не получилось считать дату. Попробуйте еще раз`,
		`response.summary.header`:   `Итого`,
		`response.task.list_header`: `*Сделать задачи:*`,
		`response.task.list_empty`:  `Нет незавершенных задач.`,
		`response.task.list_item`:   `%d. _%s_: сделаешь до "%s", получишь "%d" опыта`,

		`request.proceed`:         `*Продолжить?* %s/%s`,
		`request.task.title`:      `Введите загловок`,
		`request.task.exp`:        `Введите количество опыта полученное за выполнение задачи`,
		`request.task.expiration`: `Введите конечную дату для выполнения задачи (форматы: "%s", "%s") или напишите "%s"`,
	}

	return &msg
}
