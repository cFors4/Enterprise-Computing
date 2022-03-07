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
	REGION = "uksouth"
	URI    = "https://" + REGION + ".stt.speech.microsoft.com/" +
		"speech/recognition/conversation/cognitiveservices/v1?" +
		"language=en-US"
	KEY = "d76745e51adf4408b1f29d7a4362dc39"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func SpeechToText(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if speech, ok := t["speech"].([]byte); ok {
			if text, err := Service(speech); err == nil {
				u := map[string]interface{}{"text": text}
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

func Service(speech []byte) (string, error) {
	client := &http.Client{}
	if req, err := http.NewRequest("GET", URI, bytes.NewReader(speech)); err == nil {
		req.Header.Set("Content-Type", "audio/wav;codecs=audio/pcm;samplerate=16000")
		req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
		if rsp, err := client.Do(req); err == nil {
			defer rsp.Body.Close()
			if rsp.StatusCode == http.StatusOK {
				body, err3 := ioutil.ReadAll(rsp.Body)
				check(err3)
				return string(body), nil
			} else {
				return "", errors.New("cannot convert to speech to text")
			}
		}
	}
	return "", errors.New("Service")
}

func main() {
	r := mux.NewRouter()
	// document
	r.HandleFunc("/stt", SpeechToText).Methods("POST")
	http.ListenAndServe(":3002", r)
}
