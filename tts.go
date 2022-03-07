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
	REGIONTTS = "uksouth"
	URITTS    = "https://" + REGIONTTS + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
	KEYTTS = "d76745e51adf4408b1f29d7a4362dc39"
)

type request_struct struct {
	Text []byte `json:"text"`
}

func TextToSpeech(w http.ResponseWriter, r *http.Request) {
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
		var t request_struct
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, "Could not decode request JSON", http.StatusBadRequest)
		}
		text := t.Text
		speechOutput := ServiceTTS(text)
		fmt.Fprintf(w, speechOutput)
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func ServiceTTS(text []byte) string {
	fmt.Printf("%s\n", URITTS)
	client := &http.Client{}
	fmt.Printf("%s\n", URITTS)
	req, err := http.NewRequest("POST", URITTS, bytes.NewBuffer(text))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEYTTS)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
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
	port := 3003
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/tts", TextToSpeech)
	fmt.Printf("Listening on port: %d\n", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
