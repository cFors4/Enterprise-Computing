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
	ALPHAURI = "https:localhost:3001/alpha"
	STTURI   = "https:localhost:3002/stt"
	TTSURI   = "https:localhost:3003/tts"
)

type request struct {
	Text []byte `json:"speech"`
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
		var t request
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, "Could not decode request JSON", http.StatusBadRequest)
		}
		speech := t.Text
		if speechOutput, err := ServiceSTS(speech); err == nil {
			u := map[string]interface{}{"text": speechOutput}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(u)
			return
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func ServiceSTS(speech []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", STTURI, bytes.NewReader(speech))
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

	t := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&t); err == nil {
		return t["DisplayText"], nil
	}
	// just return DisplayText as text json
	return nil, errors.New("cannot convert speech to text")
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
