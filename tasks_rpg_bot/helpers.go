package main

import "encoding/json"

// some helper functions

// encodeToJson converts value to JSON format and logs message if fails
func encodeToJson(v interface{}) []byte {
	if data, err := json.Marshal(v); err == nil {
		return data
	} else {
		logger.Info(`Failed encoding to JSON: %s`, err)

		return data
	}
}
