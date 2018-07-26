package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// logs debug messages
const VerbosityDebug = 4

// logs normal and error messages
const VerbosityNormal = 3

// logs only error messages
const VerbosityError = 2

// logs no messages
const VerbosityQuiet = 1

// verbosity was not set
const VerbosityNotSet = 0

type LoggerConfigLog struct {
	Writer io.Writer // log writer
	Prefix string    // prefix for loader
	Flag   int       // flags for logger
}

// LoggerConfig configures settings for new logger
type LoggerConfig struct {
	Name           string // logger name
	VerbosityLevel int8   // verbosity level
	Color          bool   // colorify output

	LogStd LoggerConfigLog // log normal messages
	LogErr LoggerConfigLog // log errors
}

// Logger logs data in proper format
type Logger struct {
	logStd         *log.Logger // log normal
	logErr         *log.Logger // log errors
	Name           string
	VerbosityLevel int8
	colorizer      *ColorizerStruct
}

// DebugMode shows true when verbosity level is debug
func (l *Logger) DebugLevel() bool {
	return l.VerbosityLevel == VerbosityDebug
}

// NormalLevel shows true when verbosity level is normal
func (l *Logger) NormalLevel() bool {
	return l.VerbosityLevel == VerbosityNormal
}

// DebugMode shows true when verbosity level
func (l *Logger) ErrorLevel() bool {
	return l.VerbosityLevel == VerbosityError
}

// DebugMode shows true when verbosity level
func (l *Logger) QuietLevel() bool {
	return l.VerbosityLevel == VerbosityQuiet
}

// Info log debug message
func (l *Logger) Debug(msg string, params ...interface{}) {
	if l.VerbosityLevel < VerbosityDebug {
		return // do not log
	}

	if l.colorizer != nil {
		msg = l.colorizer.Wrap(msg, `debug`)
	}

	if len(params) > 0 {
		l.logStd.Printf(`DEBUG: `+msg+"\n", params)
	} else {
		l.logStd.Println(`DEBUG: ` + msg)
	}
}

// Info log info message
func (l *Logger) Info(msg string, params ...interface{}) {
	if l.VerbosityLevel < VerbosityNormal {
		return // do not log
	}

	if l.colorizer != nil {
		msg = l.colorizer.Wrap(msg, `info`)
	}

	if len(params) > 0 {
		l.logStd.Printf(`INFO: `+msg+"\n", params...)
	} else {
		l.logStd.Println(`INFO: ` + msg)
	}
}

// Info log error message
func (l *Logger) Error(msg string, params ...interface{}) {
	if l.VerbosityLevel < VerbosityError {
		return // do not log
	}

	if l.colorizer != nil {
		msg = l.colorizer.Wrap(msg, `error`)
	}

	if len(params) > 0 {
		l.logErr.Printf(`ERROR: `+msg+"\n", params...)
	} else {
		l.logErr.Println(`ERROR: ` + msg)
	}
}

// Info log fatal message and exit afterwards
func (l *Logger) Fatal(msg string, params ...interface{}) {
	if l.colorizer != nil {
		msg = l.colorizer.Wrap(msg, `error`)
	}

	if len(params) > 0 {
		l.logErr.Fatalf(`FATAL: `+msg+"\n", params)
	} else {
		l.logErr.Fatalln(`FATAL: ` + msg)
	}
}

// NewLogger creates new instance of logger
func NewLogger(config LoggerConfig) *Logger {
	var logger Logger

	if config.Name == `` {
		logger.Name = `defaultLogger`
	} else {
		logger.Name = config.Name
	}

	if config.Color {
		logger.colorizer = NewColorizer(map[string]string{`debug`: clGreen, `info`: clBlue, `error`: clRed})
	}

	if config.VerbosityLevel == VerbosityNotSet {
		logger.VerbosityLevel = VerbosityNormal
	} else {
		logger.VerbosityLevel = config.VerbosityLevel
	}

	if !logger.QuietLevel() {
		initLogStd(&config, &logger)
		initLogErr(&config, &logger)
	} else {
		logger.logStd = log.New(ioutil.Discard, ``, 0)
		logger.logErr = log.New(ioutil.Discard, ``, 0)
	}

	return &logger
}

func initLogErr(config *LoggerConfig, logger *Logger) {
	var errFlag int
	if config.LogErr.Writer == nil {
		config.LogErr.Writer = os.Stderr
	}
	if config.LogErr.Flag == 0 {
		errFlag = log.LstdFlags
	}
	logger.logErr = log.New(config.LogErr.Writer, config.LogErr.Prefix, errFlag)
}

func initLogStd(config *LoggerConfig, logger *Logger) {
	var stdFlag int
	if config.LogStd.Writer == nil {
		config.LogStd.Writer = os.Stdout
	}
	if config.LogStd.Flag == 0 {
		stdFlag = log.LstdFlags
	}
	logger.logStd = log.New(config.LogStd.Writer, config.LogStd.Prefix, stdFlag)
}
