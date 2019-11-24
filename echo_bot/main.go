package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

//const apiBaseUri = `http://localhost:3000/bot`
const apiBaseUri = `https://api.telegram.org/bot`
const responseTimeout = 5

var AuthKey string
var updatesOffset uint32

var botProfile struct {
	Ok     bool
	Result struct {
		Id        uint32
		FirstName string `json:"first_name"`
		Username  string
	}
}

var aliases map[string]string

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
				cmd := upd.Message.Text[ent.Offset : ent.Offset+ent.Length]
				log.Println(`Is bot command:`, cmd)

				switch cmd {
				case `/start`:
					text = `Hi, ` + upd.Message.From.FirstName + ` ` + upd.Message.From.LastName + `. Thanks for testing echo bot with us!`
				case `/time`:
					text = `*Bot time:* ` + time.Now().Format("2006-01-02 15:04:05")
				case `/code`:
					if len(upd.Message.Text) <= ent.Offset+ent.Length+1 {
						text = `*ERROR:* No input...`
					} else {
						text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
						text = "```\n" + text + "\n```"
					}
				case `/sh`:
					if len(upd.Message.Text) <= ent.Offset+ent.Length+1 {
						text = `*ERROR:* No input...`
					} else {
						query := strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])

						if aliasQuery, ok := aliases[query]; ok {
							text = execQuery(aliasQuery)
						} else {
							text = execQuery(query)
						}

					}

				case `/alias`:
					if len(upd.Message.Text) <= ent.Offset+ent.Length+1 {
						text = "*Aliases:*\n"

						if len(aliases) > 0 {
							text += "```\n"
							for k, v := range aliases {
								text += fmt.Sprintf("%s: %s\n", k, v)
							}
							text += "\n```"
						} else {
							text += `None`
						}

					} else {
						aliasKey := strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])

						if query, ok := aliases[aliasKey]; ok {
							text += "```\n"
							text += fmt.Sprintf("%s: %s\n", aliasKey, query)
							text += "\n```"
						} else {
							text = fmt.Sprintf(`*ERROR:* Unknown alias "%s"...`, aliasKey)
						}
					}

				default:
					if len(upd.Message.Text) > ent.Offset+ent.Length {
						text = strings.TrimSpace(upd.Message.Text[ent.Offset+ent.Length+1:])
					} else {
						text = `Sorry, cannot process your command`
					}
				}
				log.Println(`Message changed to:`, text)
			} else if !ent.allowedType() {
				log.Println(`-- Warning! Unexpected MessageEntity type:`, ent.Type)
			}
		}

		msg := NewSendMessage(upd.Message.Chat.Id, text /*, upd.Message.MessageId*/)

		payload := url.Values{}
		for name, value := range msg {
			payload.Set(name, value)
		}

		if _, ok := sendPostRequest(getSendMessageUrl(), []byte(payload.Encode())); !ok {
			log.Printf("Failed to send message: %s\n", payload)
		}

		sentOnceSuccessfully = true
		if upd.UpdateId >= updatesOffset {
			updatesOffset = upd.UpdateId + 1
		}
	}

	return sentOnceSuccessfully
}

func execQuery(query string) string {
	text := "```\n$ " + query + "\n"

	cmd := exec.Command(`/bin/bash`, `-c`, query)
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	
	text += "\n" + string(out)

	if err != nil {
		text += "\nERROR: " + err.Error()
	}

	text += "\n```"

	return text
}

func processRequests() {
	terminated := false

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf(`Signal "%s" called`+"\n", sig)
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

var debug bool

func main() {
	var aliasesPath string

	flag.BoolVar(&debug, `debug`, false, `Enable debug mode`)
	flag.StringVar(&aliasesPath, `aliases`, ``, `File with the list of supported aliases in JSON format as key-values`)
	flag.Parse()

	if debug {
		fmt.Println(`Started in Debug mode...`)
	} else {
		log.SetOutput(ioutil.Discard)
	}

	if len(flag.Args()) < 1 {
		log.Println(`Auth key is not defined in list of args`)
		os.Exit(1)
	}

	if aliasesPath != `` {
		_, err := os.Stat(aliasesPath)
		if os.IsNotExist(err) {
			log.Printf(`Path with aliases "%s" could not be found.`, aliasesPath)
			os.Exit(1)
		}

		aliasesContents, err := ioutil.ReadFile(aliasesPath)

		if err != nil {
			log.Printf("Error appeared during reading file: %s\n", aliasesPath)
			os.Exit(1)
		}

		if err = json.Unmarshal(aliasesContents, &aliases); err != nil {
			log.Printf("Error parsing JSON file \"%s\": %s\n", aliasesPath, err)

			os.Exit(1)
		}
	}

	AuthKey = flag.Arg(0)
	log.Println(`Echo bot started at`, time.Now().Format("2006-01-02 15:04:05"))
	log.Println(`Used auth key:`, AuthKey)
	log.Println(`Aliases path:`, aliasesPath)

	if debug {
		log.Println(`Aliases:`)

		for k, v := range aliases {
			log.Println(` `, k, "=", v)
		}
	}

	if checkConnection() {
		processRequests()
	}

	log.Println(`Echo bot finished at`, time.Now().Format("2006-01-02 15:04:05"))
}
