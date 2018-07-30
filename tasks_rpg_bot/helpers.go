package main

import "encoding/json"

// some helper functions

// encodeToJson converts value to JSON format and logs message if fails
func encodeToJson(v interface{}) string {
	if data, err := json.Marshal(v); err == nil {
		return string(data)
	} else {
		logger.Info(`Failed encoding to JSON: %s`, err)

		return string(data)
	}
}
