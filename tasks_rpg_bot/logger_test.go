package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestLoggerName(t *testing.T) {
	loggerName := `testLog`
	errOutput := bytes.NewBufferString(``)
	stdOutput := bytes.NewBufferString(``)
	loggerCfg := LoggerConfig{
		Name: loggerName,
		LogStd: LoggerConfigLog{
			Writer: stdOutput,
		},
		LogErr: LoggerConfigLog{
			Writer: errOutput,
		},
	}
	logger := NewLogger(loggerCfg)

	if logger.Name != loggerName {
		t.Errorf("Expecting logger name set to \"%s\", got \"%v\"\n", loggerName, logger.Name)
	}
}

func TestLoggerOutputDefaults(t *testing.T) {
	errOutput := bytes.NewBufferString(``)
	stdOutput := bytes.NewBufferString(``)

	loggerCfg := LoggerConfig{
		LogStd: LoggerConfigLog{
			Writer: stdOutput,
		},
		LogErr: LoggerConfigLog{
			Writer: errOutput,
		},
	}
	logger := NewLogger(loggerCfg)

	logger.Info(`Good job, %s!`, `Man`)
	logger.Info(`Nice work!`)

	if !strings.Contains(stdOutput.String(), `Good job, Man!`) ||
		!strings.Contains(stdOutput.String(), `Nice work!`) {
		t.Error(`Standard output does not contain expected log messages`)
	}

	logger.Error(`Keep doing, %s!`, `Man`)
	logger.Error(`Nicely done!`)

	if !strings.Contains(errOutput.String(), `Keep doing, Man!`) ||
		!strings.Contains(errOutput.String(), `Nicely done!`) {
		t.Error(`Error output does not contain expected log messages`)
	}
}
