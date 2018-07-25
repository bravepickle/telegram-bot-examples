package tasks_rpg_bot

import (
	"io"
	"io/ioutil"
	"log"
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

	LogStd LoggerConfigLog // log normal messages
	LogErr LoggerConfigLog // log errors
}

// Logger logs data in proper format
type Logger struct {
	logStd         *log.Logger // log normal
	logErr         *log.Logger // log errors
	Name           string
	VerbosityLevel int8
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

	if len(params) > 0 {
		l.logStd.Printf(msg+"\n", params)
	} else {
		l.logStd.Println(msg)
	}
}

// Info log info message
func (l *Logger) Info(msg string, params ...interface{}) {
	if l.VerbosityLevel < VerbosityNormal {
		return // do not log
	}

	if len(params) > 0 {
		l.logStd.Printf(msg+"\n", params...)
	} else {
		l.logStd.Println(msg)
	}
}

// Info log error message
func (l *Logger) Error(msg string, params ...interface{}) {
	if l.VerbosityLevel < VerbosityError {
		return // do not log
	}

	if len(params) > 0 {
		l.logErr.Printf(msg+"\n", params...)
	} else {
		l.logErr.Println(msg)
	}
}

// Info log fatal message and exit afterwards
func (l *Logger) Fatal(msg string, params ...interface{}) {
	if len(params) > 0 {
		l.logErr.Fatalf(msg+"\n", params)
	} else {
		l.logErr.Fatalln(msg)
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
		panic(`Error log output not defined`)
	}
	if config.LogErr.Flag == 0 {
		errFlag = log.LstdFlags
	}
	logger.logErr = log.New(config.LogErr.Writer, config.LogErr.Prefix, errFlag)
}

func initLogStd(config *LoggerConfig, logger *Logger) {
	var stdFlag int
	if config.LogStd.Writer == nil {
		panic(`Standard log output not defined`)
	}
	if config.LogStd.Flag == 0 {
		stdFlag = log.LstdFlags
	}
	logger.logStd = log.New(config.LogStd.Writer, config.LogStd.Prefix, stdFlag)
}