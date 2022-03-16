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
		if textOutput, err := handleSTT(speech); err == nil {
			if answerOutput, err := handleALPHA(textOutput); err == nil {
				u := map[string]interface{}{"text": answerOutput}
				// if speechOutput, err := handleTTS(answerOutput); err == nil {
				// 	w.WriteHeader(http.StatusOK)
				// 	json.NewEncoder(w).Encode(speechOutput)
				// 	return
				// }
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(u)
				return
			}
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported")
	}
}

func handleSTT(speech []byte) (string, error) {
	body := &request_speech{
		Text: speech,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", STTURI, payloadBuf)
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
		fmt.Printf(string(body))
		parts := strings.Split(string(body), ":")
		replaceValue := strings.Replace(parts[1], "}", "", -1)
		retrunValue := strings.Trim(replaceValue, "\"")
		retrunValue2 := strings.Replace(retrunValue, "\"", "", -1)
		retrunValue3 := strings.Replace(retrunValue2, "\n", "", -1)
		//handle output to remove json formatting
		return retrunValue3, nil
	}

	fmt.Printf("Cannot convert Question to Anwer")
	return "", errors.New("Cannot convert Question to Anwer")
}

func handleALPHA(text string) (string, error) {
	body := &request_text{
		Text: text,
	}
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", ALPHAURI, payloadBuf)
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
		return string(body), nil
	}

	fmt.Printf("Cannot convert Question to Anwer")
	return "", errors.New("Cannot convert Question to Anwer")
}

// func handleTTS(speech string) (string, error) {
// 	client := &http.Client{}
// 	req, err := http.NewRequest("POST", TTSURI, bytes.NewReader(speech))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", errors.New("HTTP Client cannot carry out request")
// 	}
// 	defer resp.Body.Close()

// 	if body, err := ioutil.ReadAll(resp.Body); err == nil {
// 		return string(body), nil
// 	}

// 	return "", errors.New("cannot convert Question to Anwer")
// }

func main() {
	port := 3000
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/alexa", SpeechToSpeech)
	fmt.Printf("Listening on port: %d\n", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
