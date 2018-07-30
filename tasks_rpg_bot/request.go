package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type RequestManagerStruct struct{}

func (r *RequestManagerStruct) SendGetRequest(uri string) ([]byte, bool) {
	logger.Info(`Calling GET %s`, uri)

	resp, err := http.Get(uri)

	var body []byte

	if err != nil {
		logger.Error("Error response message: %s", err)

		return body, false
	}

	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Response returned non-OK status code (%d) with body: %s", resp.StatusCode, body)

		return body, false
	}

	if err != nil {
		logger.Error("Error reading response body: %s", err)

		return body, false
	}

	logger.Debug("Response: %s", body)

	return body, true
}

func (r *RequestManagerStruct) SendPostJsonRequest(uri string, data interface{}) ([]byte, bool) {
	payload := encodeToJson(data)

	logger.Info(`Calling JSON POST %s with payload: %s`, uri, string(payload))

	resp, err := http.Post(uri, `application/json`, bytes.NewBuffer(payload))

	var body []byte

	if err != nil {
		logger.Error("Error response message: %s", err)

		return body, false
	}

	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Response returned nok-OK status code (%d) with body: %s", resp.StatusCode, body)

		return body, false
	}

	if err != nil {
		logger.Error("Error reading response body: %s", err)

		return body, false
	}

	logger.Debug("Response: %s", body)

	return body, true
}

func (r *RequestManagerStruct) SendPostRequest(uri string, payload []byte) ([]byte, bool) {
	logger.Info(`Calling POST %s with %s`, uri, string(payload))

	resp, err := http.Post(uri, `application/x-www-form-urlencoded`, bytes.NewBuffer(payload))

	var body []byte

	if err != nil {
		logger.Error("Error response message: %s", err)

		return body, false
	}

	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Response returned nok-OK status code (%d) with body: %s", resp.StatusCode, body)

		return body, false
	}

	if err != nil {
		logger.Error("Error reading response body: %s", err)

		return body, false
	}

	logger.Debug("Response: %s", body)

	return body, true
}
