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
	Text string `json:"text"`
}

func TextToSpeech(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tts" {
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
		byteText := []byte(text)
		speechOutput := ServiceTTS(byteText)
		err3 := ioutil.WriteFile("speech.wav", speechOutput, 0644)
		if err3 != nil {
			http.Error(w, "Could not write speech file", http.StatusBadRequest)
		}
		fmt.Println(speechOutput)
		fmt.Fprintf(w, string(speechOutput))
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func ServiceTTS(text []byte) []byte {
	client := &http.Client{}
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
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return body
	}

	return nil
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
