package main

import (
	"bytes"
	"encoding/json"
)

// some helper functions

// encodeToJson converts value to JSON format and logs message if fails
func encodeToJson(v interface{}) []byte {
	if data, err := json.Marshal(v); err == nil {
		if appConfig.GetAppJsonPretty() {
			var buffer bytes.Buffer
			if err = json.Indent(&buffer, data, ``, "  "); err != nil {
				logger.Error(`Failed to pretty-print JSON: %s`, err)
			}

			return buffer.Bytes()
		}

		return data
	} else {
		logger.Info(`Failed encoding to JSON: %s`, err)

		return data
	}
}
