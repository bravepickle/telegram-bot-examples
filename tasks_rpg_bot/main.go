package main

import (
	"flag"
	"fmt"
	"os"
)

var logger *Logger
var inputFlags inputFlagsStruct
var api *TelegramBotsApiStruct
var appConfig *AppConfigStruct

const pollSleepInterval = 5

func init() {
	flag.BoolVar(&inputFlags.Debug, "v", false, "Verbose output.")
	flag.BoolVar(&inputFlags.Error, "e", false, "Errors output only.")
	flag.BoolVar(&inputFlags.Quiet, "q", false, "No output.")
	flag.BoolVar(&inputFlags.Color, "c", false, "Disable colored output.")
	flag.StringVar(&inputFlags.AuthKey, "k", ``, "Telegram Bots API Auth Key. Required.")
	flag.IntVar(&inputFlags.Sleep, "s", pollSleepInterval, "Sleep interval in seconds between polling for updates.")
}

// TODO: init logs and start bot run based on configs

func main() {
	flag.Parse()

	fmt.Printf("Starting application in \"%s\" mode\n", inputFlags.StringVerbosity())
	//fmt.Printf("Color mode is \"%t\"\n", inputFlags.Color)

	appConfig = NewAppConfig()

	logger = initLogger()

	api = NewTelegramBotsApi(getAuthKey(), inputFlags.Sleep)

	if !api.checkConnection() {
		logger.Fatal(`Failed to establish connection to %s`, api)
	}

	logger.Info(`Successfully connected to %s`, api.BotInfo.Result.Username)

	api.processRequests()

	////logger.Debug(`API ME: ` + api.routingMe.Uri())
	//logger.Debug(`API ME: ` + api.routingMe.Uri())
	//logger.Debug(`Hello, weak world!`)
	//logger.Info(`Hello, world!`)
	//logger.Error(`Hello, bad world!`)
	////logger.Fatal(`Hello, chaos world! %s`, *api)
	//logger.Fatal(`Hello, chaos world!`)

}

func getAuthKey() string {
	var authKey string
	if inputFlags.AuthKey == `` {
		authKey = appConfig.GetApiAuthKey()

		if authKey == `` {
			fmt.Println(logger.colorizer.Wrap(`Error! Auth Key value not set neither in config nor in command options`, `error`))
			flag.Usage()
			os.Exit(1)
		}
	} else {
		authKey = inputFlags.AuthKey
	}

	return authKey
}

func initLogger() *Logger {
	verbosityLevel := inputFlags.ParseVerbosity()

	logger := NewLogger(LoggerConfig{
		VerbosityLevel: verbosityLevel,
		Color:          inputFlags.Color,

		LogErr: LoggerConfigLog{
		//Prefix: `ERR: `,
		}, LogStd: LoggerConfigLog{
		//Prefix: `OUT: `,
		},
	})

	return logger
}
