package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const cfgAppLocale = `APP_LOCALE`
const cfgAppJsonPretty = `APP_JSON_PRETTY`
const cfgDbDsn = `DB_DSN`
const cfgApiAuthKey = `API_AUTH_KEY`
const cfgApiBaseUri = `API_BASE_URI`
const cfgApiTimeout = `API_TIMEOUT`
const cfgApiUpdatesInterval = `API_UPDATES_INTERVAL`
const cfgApiRemindInterval = `API_REMIND_INTERVAL`

const defaultEnvFile = `.env`
const defaultApiRemindIntervalHr = 24   // in hours
const defaultApiUpdatesIntervalSec = 10 // in seconds
const defaultResponseTimeout = 5        // in seconds
const defaultApiBaseUri = `https://api.telegram.org/bot`

type AppConfigStruct struct {
	params map[string]string
}

func (c *AppConfigStruct) Get(param string, defaultVal string) string {
	if foundVal, ok := c.params[param]; ok {
		return foundVal
	}

	return defaultVal
}

// GetAppLocale get DB DSN
func (c *AppConfigStruct) GetAppLocale() string {
	return c.Get(cfgAppLocale, defaultLocale)
}

// GetAppJsonPretty get JSON pretty print flag
func (c *AppConfigStruct) GetAppJsonPretty() bool {
	return c.getBoolValue(cfgAppJsonPretty, false)
}

// GetDbDsn get DB DSN
func (c *AppConfigStruct) GetDbDsn() string {
	return c.Get(cfgDbDsn, ``)
}

// GetApiBaseUri get Telegram Bot API base URI
func (c *AppConfigStruct) GetApiBaseUri() string {
	return c.Get(cfgApiBaseUri, defaultApiBaseUri)
}

// GetApiTimeout get Telegram Bot API secret
func (c *AppConfigStruct) GetApiAuthKey() string {
	return c.Get(cfgApiAuthKey, ``)
}

// GetApiTimeout get API response timeout
func (c *AppConfigStruct) GetApiTimeout() int {
	return c.getIntValue(cfgApiTimeout, defaultResponseTimeout)
}

// GetApiRemindInterval returns number of seconds interval between reminding on unfinished tasks
func (c *AppConfigStruct) GetApiUpdatesInterval() int {
	return c.getIntValue(cfgApiUpdatesInterval, defaultApiUpdatesIntervalSec)
}

// GetApiRemindInterval returns number of hours interval between reminding on unfinished tasks
func (c *AppConfigStruct) GetApiRemindInterval() int {
	return c.getIntValue(cfgApiRemindInterval, defaultApiRemindIntervalHr)
}

// load loads env files to struct
func (c *AppConfigStruct) load(filenames ...string) {
	data, err := godotenv.Read(filenames...)
	if err != nil {
		log.Fatalf(`Failed to load config(s): %s`, err)
	}

	c.params = data
}

// getIntValue get value of type integer routine
func (c *AppConfigStruct) getIntValue(name string, defaultValue int) int {
	rawValue := c.Get(name, ``)
	if rawValue == `` {
		return defaultValue
	}

	value, err := strconv.Atoi(rawValue)
	if err != nil {
		logger.Fatal(`Failed to parse rawValue for "%s": %s`, name, err)
	}

	return value
}

// getBoolValue get value of type bool routine
func (c *AppConfigStruct) getBoolValue(name string, defaultValue bool) bool {
	rawValue := c.Get(name, ``)
	if rawValue == `` {
		return defaultValue
	}

	if rawValue == `0` {
		return false
	}

	if rawValue == `1` {
		return true
	}

	return defaultValue
}

func NewAppConfig() *AppConfigStruct {
	var config AppConfigStruct

	filename := defaultEnvFile

	if _, err := os.Stat(filename); err != nil {
		log.Fatalf(`Failed to read file "%s" with error: %s`, filename, err)
	}

	config.load(filename)

	return &config
}
