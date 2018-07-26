package main

import (
	"flag"
	"fmt"
	"os"
)

var logger *Logger
var inputFlags inputFlagsStruct
var api *TelegramBotsApiStruct

func init() {
	flag.BoolVar(&inputFlags.Debug, "v", false, "Verbose output.")
	flag.BoolVar(&inputFlags.Error, "e", false, "Errors output only.")
	flag.BoolVar(&inputFlags.Quiet, "q", false, "No output.")
	flag.BoolVar(&inputFlags.Color, "c", false, "Disable colored output.")
	flag.StringVar(&inputFlags.AuthKey, "k", ``, "Telegram Bots API Auth Key. Required.")
}

// TODO: init logs and start bot run based on configs

func main() {
	flag.Parse()

	fmt.Printf("Starting application in \"%s\" mode\n", inputFlags.StringVerbosity())
	//fmt.Printf("Color mode is \"%t\"\n", inputFlags.Color)

	logger = initLogger()

	if inputFlags.AuthKey == `` {
		fmt.Println(logger.colorizer.Wrap(`Error! Required option Auth Key value not set`, `error`))
		flag.Usage()
		os.Exit(1)
	}

	api = NewTelegramBotsApi(inputFlags.AuthKey)

	if !api.checkConnection() {
		logger.Fatal(`Failed to establish connection to %s`, api)
	}

	logger.Info(`Successfully connected to %s`, api.BotInfo.Result.Username)

	////logger.Debug(`API ME: ` + api.routingMe.Uri())
	//logger.Debug(`API ME: ` + api.routingMe.Uri())
	//logger.Debug(`Hello, weak world!`)
	//logger.Info(`Hello, world!`)
	//logger.Error(`Hello, bad world!`)
	////logger.Fatal(`Hello, chaos world! %s`, *api)
	//logger.Fatal(`Hello, chaos world!`)

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
