package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	ALPHAURI = "http://localhost:3001/alpha"
	STTURI   = "http://localhost:3002/stt"
	TTSURI   = "http://localhost:3003/tts"
)

type request_speech struct {
	Text []byte `json:"speech"`
}

type request_text struct {
	Text string `json:"text"`
}

func SpeechToSpeech(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/alexa" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not decode request body", http.StatusBadRequest)
			return
		}
		var t request_speech
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, "Could not decode request JSON", http.StatusBadRequest)
		}
		speech := t.Text
		if textOutput, err := handleService(STTURI, handleSpeech(speech)); err == nil {
			if answerOutput, err := handleService(ALPHAURI, handleText(textOutput)); err == nil {
				if speechOutput, err := handleService(TTSURI, handleText(answerOutput)); err == nil {
					fmt.Printf(speechOutput)
					u := map[string]interface{}{"speech": speechOutput}
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(u)
					return
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func handleText(text string) *bytes.Buffer {
	body := &request_text{
		Text: text,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	return payloadBuf

}
func handleSpeech(speech []byte) *bytes.Buffer {
	body := &request_speech{
		Text: speech,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	return payloadBuf
}

func handleService(URI string, payloadBuf *bytes.Buffer) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", URI, payloadBuf)
	if err != nil {
		fmt.Printf("Cannot display request")
		return "", errors.New("Cannot display request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP Client cannot carry out request")
		return "", errors.New("HTTP Client cannot carry out request stt")
	}
	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		returnValue := JSONtoStringValue(body)
		fmt.Printf("\n" + URI + returnValue)
		return returnValue, nil
	}
	fmt.Printf("Cannot convert from endpoint")
	return "", errors.New("Cannot convert Speech to Text")
}

func JSONtoStringValue(body []byte) string {
	parts := strings.Split(string(body), ":")
	replacePartValue := strings.Replace(parts[1], "}", "", -1)
	trimValue := strings.Trim(replacePartValue, "\"")
	replaceValue := strings.Replace(trimValue, "\"", "", -1)
	retrunValue := strings.Replace(replaceValue, "\n", "", -1)
	return retrunValue
}

func main() {
	port := 3000
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/alexa", SpeechToSpeech)
	fmt.Printf("Listening on port: %d\n", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
