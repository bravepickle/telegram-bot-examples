package main

import (
	"flag"
	"fmt"
)

var logger *Logger
var inputFlags inputFlagsStruct

func init() {
	flag.BoolVar(&inputFlags.Debug, "v", false, "Verbose output.")
	flag.BoolVar(&inputFlags.Error, "e", false, "Errors output only.")
	flag.BoolVar(&inputFlags.Quiet, "q", false, "No output.")
	flag.BoolVar(&inputFlags.Color, "c", false, "Disable colored output.")
}

// TODO: init logs and start bot run based on configs

func main() {
	flag.Parse()

	fmt.Printf("Starting application in \"%s\" mode\n", inputFlags.StringVerbosity())
	//fmt.Printf("Color mode is \"%t\"\n", inputFlags.Color)

	logger = initLogger()

	logger.Debug(`Hello, weak world!`)
	logger.Info(`Hello, world!`)
	logger.Error(`Hello, bad world!`)
	logger.Fatal(`Hello, chaos world!`)
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
