package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

const (
	REGIONTTS = "uksouth"
	URITTS    = "https://" + REGIONTTS + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
	KEYTTS = "d76745e51adf4408b1f29d7a4362dc39"
)

type request_T struct {
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
		var t request_T
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, "Could not decode request JSON", http.StatusBadRequest)
		}
		text := t.Text
		if speech, err := ServiceTTS(text); err == nil {
			u := map[string]interface{}{"speech": speech}
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
func ServiceTTS(input string) ([]byte, error) {
	client := &http.Client{}
	createXMLInput(input)
	text, err := ioutil.ReadFile("text3.xml")

	req, err := http.NewRequest("POST", URITTS, bytes.NewBuffer(text))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEYTTS)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
	//DSIPLAY request server-side
	res, err := httputil.DumpRequest(req, true)
	if err != nil {
		return nil, errors.New("Cannot display request")
	}
	fmt.Println("Sent Request:")
	fmt.Print(string(res))

	rsp, err2 := client.Do(req)
	if err2 != nil {
		return nil, errors.New("HTTP Client cannot carry out request")
	}

	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		if err3 != nil {
			return nil, errors.New("Cannot read response")
		}
		ioutil.WriteFile("speech.wav", body, 0644)
		return body, nil
	} else {
		fmt.Printf("\nStatus Code: \n %d\n", rsp.StatusCode)
		return nil, errors.New("Cannot convert text to speech")
	}
}

func createXMLInput(input string) {
	f, err := os.Create("text3.xml")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	d := []string{"<?xml version='1.0'?>",
		"<speak version='1.0' xml:lang='en-US'>",
		"<voice xml:lang='en-US' name='en-US-JennyNeural'>",
		input,
		"</voice>",
		"</speak>"}

	for _, v := range d {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file written successfully")
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
