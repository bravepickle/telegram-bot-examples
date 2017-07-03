package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
	"strings"
)

//const apiBaseUri = `http://localhost:3000/bot`
const apiBaseUri = `https://api.telegram.org/bot`
const responseTimeout = 5

var AuthKey string
var updatesOffset uint32

var botProfile struct {
	Ok bool
	Result struct {
		Id        uint32
		FirstName string `json:"first_name"`
		Username  string
	}
}

func getUpdatesUrl() string {
	return apiBaseUri + AuthKey + `/getUpdates?timeout=` + strconv.Itoa(responseTimeout) + `&offset=` + strconv.FormatInt(int64(updatesOffset), 10)
}

func getMeUrl() string {
	return apiBaseUri + AuthKey + `/getMe`
}

func getSendMessageUrl() string {
	//return `http://localhost:3000/bot` + AuthKey + `/sendMessage`
	return apiBaseUri + AuthKey + `/sendMessage`
}

func readBodyFromGetRequest(uri string) ([]byte, bool) {
	log.Println(`Calling`, uri)
	resp, err := http.Get(uri)

	var body []byte

	if err != nil {
		log.Printf("Error response message: %s\n", err)

		return body, false
	}

	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Response returned nok-OK status code (%d) with body: %s\n", resp.StatusCode, body)

		return body, false
	}

	if err != nil {
		log.Printf("Error reading response body: %s\n", err)

		return body, false
	}

	log.Printf("Response: %s\n", body)

	return body, true
}

func sendPostRequest(uri string, payload []byte) ([]byte, bool) {

	log.Println(`Calling`, uri, `with`, string(payload))
	resp, err := http.Post(uri, `application/x-www-form-urlencoded`, bytes.NewBuffer(payload))

	var body []byte

	if err != nil {
		log.Printf("Error response message: %s\n", err)

		return body, false
	}

	body, err = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("Response returned nok-OK status code (%d) with body: %s\n", resp.StatusCode, body)

		return body, false
	}

	if err != nil {
		log.Printf("Error reading response body: %s\n", err)

		return body, false
	}

	log.Printf("Response: %s\n", body)

	return body, true
}

func checkConnection() bool {
	body, ok := readBodyFromGetRequest(getMeUrl())
	if !ok {
		return false
	}

	if err := json.Unmarshal(body, &botProfile); err != nil {
		log.Printf("Error parsing JSON: %s\n", err)

		return false
	}

	return true
}

func processUpdates() bool {
	body, ok := readBodyFromGetRequest(getUpdatesUrl())
	if !ok {
		return false
	}

	var updates updatesPayloadStruct

	if err := json.Unmarshal(body, &updates); err != nil {
		log.Printf("Error parsing JSON: %s\n", err)

		return false
	}

	// TODO: edited_message handle, inline_query

	sentOnceSuccessfully := false

	for _, upd := range updates.Result {
		log.Println(`Handling update `, upd.UpdateId, `message`, upd.Message.MessageId)
		log.Println(`>`, upd.Message.Text)

		var text = upd.Message.Text
		for _, ent := range upd.Message.Entities {
			if ent.Type == `bot_command` {
				cmd := upd.Message.Text[ent.Offset:ent.Offset + ent.Length]
				log.Println(`Is bot command:`, cmd)

				switch cmd {
					case `/time`:
						text = `*Bot time:* ` + time.Now().Format("2006-01-02 15:04:05")
					case `/code`:
						text = strings.TrimSpace(upd.Message.Text[ent.Offset + ent.Length + 1:])
						text = "```\n" + text + "\n```"
					default:
						text = strings.TrimSpace(upd.Message.Text[ent.Offset + ent.Length + 1:])
				}
				log.Println(`Message changed to:`, text)
			} else if !ent.allowedType() {
				log.Println(`-- Warning! Unexpected MessageEntity type:`, ent.Type)
			}
		}

		msg := NewSendMessage(upd.Message.Chat.Id, text/*, upd.Message.MessageId*/)

		payload := url.Values{}
		for name, value := range msg {
			payload.Set(name, value)
		}

		if _, ok := sendPostRequest(getSendMessageUrl(), []byte(payload.Encode())); !ok {
			log.Printf("Failed to send message: %s\n", payload)
		}

		sentOnceSuccessfully = true

		if upd.UpdateId >= updatesOffset {
			log.Printf(" --------- Was offset %d, will be: %d\n", updatesOffset, upd.UpdateId+1)
			updatesOffset = upd.UpdateId + 1

		}
	}

	return sentOnceSuccessfully
}

func processRequests() {
	terminated := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf(`Signal "%s" called`+"\n", sig)
			terminated = true
		}
	}()

	for {
		processUpdates()

		if terminated {
			break
		}

		time.Sleep(time.Second * 1)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println(`Auth key is not defined in list of args`)
		os.Exit(1)
	}

	AuthKey = os.Args[1]
	fmt.Println(`Echo bot started at`, time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(`Used auth key:`, AuthKey)

	if checkConnection() {
		processRequests()
	}

	fmt.Println(`Echo bot finished at`, time.Now().Format("2006-01-02 15:04:05"))
}
