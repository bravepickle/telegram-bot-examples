package main

import "encoding/json"

// Telegram Bots API models

type BotProfileTelegramModel struct {
	Ok     bool
	Result struct {
		Id        uint32
		FirstName string `json:"first_name"`
		Username  string
	}
}

func (p BotProfileTelegramModel) String() string {
	if text, err := json.Marshal(p); err == nil {
		return string(text)
	} else {
		logger.Error(`Error ocurred when encoding to JSON bot profile object: %s`, err)

		return `[UNDEFINED]`
	}
}
