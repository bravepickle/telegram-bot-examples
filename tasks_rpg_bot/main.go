package main

//import "telegram_bots/tasks_rpg_bot"

var logger Logger

// TODO: init logs and start bot run based on configs

func main() {
	logger := NewLogger(LoggerConfig{
		VerbosityLevel: VerbosityDebug,
		LogErr: LoggerConfigLog{
			Prefix: `ERR: `,
		}, LogStd: LoggerConfigLog{
			Prefix: `OUT: `,
		},
	})

	logger.Debug(`Hello, weak world!`)
	logger.Info(`Hello, world!`)
	logger.Error(`Hello, bad world!`)
	logger.Fatal(`Hello, chaos world!`)
}
