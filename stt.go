package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
	KEY = "d76745e51adf4408b1f29d7a4362dc39"
)

type request struct {
	Text []byte `json:"speech"`
}

func SpeechToText(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/stt" {
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
		var t request
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, "Could not decode request JSON", http.StatusBadRequest)
		}
		speech := t.Text
		textOutput := ServiceSTT(speech)
		fmt.Fprintf(w, textOutput)
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func ServiceSTT(speech []byte) string {
	client := &http.Client{}
	fmt.Printf("%s\n", URI)
	req, err := http.NewRequest("POST", URI, bytes.NewReader(speech))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type",
		"audio/wav;codecs=audio/pcm;samplerate=16000")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

func main() {
	port := 3002
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/stt", SpeechToText)
	fmt.Printf("Listening on port: %d\n", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
