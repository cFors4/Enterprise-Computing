package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	REGIONTTS = "uksouth"
	uriTTS    = "https://" + REGIONTTS + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
	KEYTTS = "d76745e51adf4408b1f29d7a4362dc39"
)

func checkTTS(e error) {
	if e != nil {
		panic(e)
	}
}

func TextToSpeech(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if text, ok := t["text"].([]byte); ok {
			if speech, err := ServiceTTS(text); err == nil {
				u := map[string]interface{}{"speech": speech}
				err3 := ioutil.WriteFile("speech.wav", speech, 0644)
				checkTTS(err3)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(u)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func ServiceTTS(text []byte) ([]byte, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("POST", uriTTS, bytes.NewBuffer(text)); err == nil {
		req.Header.Set("Content-Type", "application/ssml+xml")
		req.Header.Set("Ocp-Apim-Subscription-Key", KEYTTS)
		req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
		if rsp, err := client.Do(req); err == nil {
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				body, err3 := ioutil.ReadAll(rsp.Body)
				checkTTS(err3)
				return body, nil
			} else {
				return nil, errors.New("cannot convert to speech to text")
			}
		}
	}
	return nil, errors.New("Service")
}

func main() {
	r := mux.NewRouter()
	// document
	r.HandleFunc("/tts", TextToSpeech).Methods("POST")
	http.ListenAndServe(":3003", r)
}
