package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

const defaultEnvFile = `.env`
const cfgDbDsn = `DB_DSN`
const cfgApiAuthKey = `API_AUTH_KEY`
const cfgApiTimeout = `API_TIMEOUT`

type AppConfigStruct struct {
	params map[string]string
}

func (c *AppConfigStruct) Get(param string, defaultVal string) string {
	if foundVal, ok := c.params[param]; ok {
		return foundVal
	}

	return defaultVal
}

func (c *AppConfigStruct) GetDbDsn() string {
	return c.Get(cfgDbDsn, ``)
}

func (c *AppConfigStruct) GetApiAuthKey() string {
	return c.Get(cfgApiAuthKey, ``)
}
func (c *AppConfigStruct) GetApiTimeout() int {
	strTimeout := c.Get(cfgApiTimeout, ``)
	if strTimeout == `` {
		return responseTimeoutDefault
	}

	timeout, err := strconv.Atoi(strTimeout)

	if err != nil {
		logger.Fatal(`Failed to parse value for "%s": %s`, cfgApiTimeout, err)
	}

	return timeout
}

func (c *AppConfigStruct) load(filenames ...string) {
	data, err := godotenv.Read(filenames...)
	if err != nil {
		log.Fatalf(`Failed to load config(s): %s`, err)
	}

	//log.Printf(`CONFIG: %v` + "\n", data)

	c.params = data
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

// rootPath checks the current path
//func rootPath() string {
//	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
//	if err != nil {
//		log.Fatal(err)
//	}
//}
