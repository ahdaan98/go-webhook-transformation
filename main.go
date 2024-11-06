
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

const (
	port              = ":8080"
	webhookURL        = "https://webhook.site/0376d9ec-e76d-4aa7-9a94-ea7f32dab4ef"
	contentTypeJSON   = "application/json"
	maxEventChannel   = 100
)

type InputEvent struct {
	Event           string            `json:"ev"`
	EventType       string            `json:"et"`
	AppId           string            `json:"id"`
	UserId          string            `json:"uid"`
	MessageId       string            `json:"mid"`
	PageTitle       string            `json:"t"`
	PageURL         string            `json:"p"`
	BrowserLanguage string            `json:"l"`
	ScreenSize      string            `json:"sc"`
	DynamicFields   map[string]string `json:"-"`
}

type OutputEvent struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppId           string               `json:"app_id"`
	UserId          string               `json:"user_id"`
	MessageId       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"`
	Traits          map[string]Attribute `json:"traits"`
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

var eventChannel = make(chan InputEvent, maxEventChannel)
var sendChannel = make(chan []byte, maxEventChannel)

func main() {
	go worker()
	http.HandleFunc("/", handleRequest)

	log.Printf("Server listening on %v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Method should be POST", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Failed to read the body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	var input InputEvent
	if err := json.Unmarshal(bodyBytes, &input); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	input.DynamicFields = make(map[string]string)
	if err := json.Unmarshal(bodyBytes, &input.DynamicFields); err != nil {
		http.Error(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	eventChannel <- input

	var resp OutputEvent
	select {
	case val := <-sendChannel:
		if err := json.Unmarshal(val, &resp); err != nil {
			http.Error(res, "Failed to unmarshal response", http.StatusBadRequest)
			return
		}
	default:
		http.Error(res, "No response received", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", contentTypeJSON)
	response := map[string]interface{}{
		"message": "Request sent to webhook",
		"data":    resp,
	}

	if err := json.NewEncoder(res).Encode(response); err != nil {
		http.Error(res, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}

func worker() {
	for event := range eventChannel {
		processEvent(event)
	}
}

func processEvent(input InputEvent) {
	output := OutputEvent{
		Event:           input.Event,
		EventType:       input.EventType,
		AppId:           input.AppId,
		UserId:          input.UserId,
		MessageId:       input.MessageId,
		PageTitle:       input.PageTitle,
		PageURL:         input.PageURL,
		BrowserLanguage: input.BrowserLanguage,
		ScreenSize:      input.ScreenSize,
		Attributes:      make(map[string]Attribute),
		Traits:          make(map[string]Attribute),
	}

	attrKeys := make(map[int]string)
	attrValues := make(map[int]string)
	attrTypes := make(map[int]string)
	traitKeys := make(map[int]string)
	traitValues := make(map[int]string)
	traitTypes := make(map[int]string)

	for key, value := range input.DynamicFields {
		switch {
		case regexp.MustCompile(`^atrk(\d+)$`).MatchString(key):
			index, _ := strconv.Atoi(key[4:])
			attrKeys[index] = value
		case regexp.MustCompile(`^atrv(\d+)$`).MatchString(key):
			index, _ := strconv.Atoi(key[4:])
			attrValues[index] = value
		case regexp.MustCompile(`^atrt(\d+)$`).MatchString(key):
			index, _ := strconv.Atoi(key[4:])
			attrTypes[index] = value
		case regexp.MustCompile(`^uatrk(\d+)$`).MatchString(key):
			index, _ := strconv.Atoi(key[5:])
			traitKeys[index] = value
		case regexp.MustCompile(`^uatrv(\d+)$`).MatchString(key):
			index, _ := strconv.Atoi(key[5:])
			traitValues[index] = value
		case regexp.MustCompile(`^uatrt(\d+)$`).MatchString(key):
			index, _ := strconv.Atoi(key[5:])
			traitTypes[index] = value
		}
	}

	for index, key := range attrKeys {
		if value, ok := attrValues[index]; ok {
			if attrType, ok := attrTypes[index]; ok {
				output.Attributes[key] = Attribute{
					Value: value,
					Type:  attrType,
				}
			}
		}
	}

	for index, key := range traitKeys {
		if value, ok := traitValues[index]; ok {
			if traitType, ok := traitTypes[index]; ok {
				output.Traits[key] = Attribute{
					Value: value,
					Type:  traitType,
				}
			}
		}
	}

	outputJSON, err := json.Marshal(output)
	if err != nil {
		log.Printf("Error marshaling output: %v", err)
		return
	}

	sendChannel <- outputJSON

	resp, err := http.Post(webhookURL, contentTypeJSON, bytes.NewBuffer(outputJSON))
	if err != nil {
		log.Printf("Error sending to webhook: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Webhook returned non-200 status: %d", resp.StatusCode)
	}
}