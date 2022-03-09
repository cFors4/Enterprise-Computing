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
)

const (
	REGIONTTS = "uksouth"
	URITTS    = "https://" + REGIONTTS + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
	KEYTTS = "d76745e51adf4408b1f29d7a4362dc39"
)

type request struct {
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
		// decode json
		var t request
		err = json.Unmarshal(body, &t)
		if err != nil {
			http.Error(w, "Could not decode request JSON", http.StatusBadRequest)
		}
		// send request
		text := t.Text
		byteText := []byte(text)
		if speech, err := ServiceTTS(byteText); err == nil {
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

//Handle request TTS
func ServiceTTS(text []byte) (interface{}, error) {
	//Create POST request with URI, text input, and headers
	client := &http.Client{}
	fmt.Printf("%s\n", URITTS)
	//TODO: Modify text.xml to have text parameter
	text, err := ioutil.ReadFile("text.xml")
	req, err := http.NewRequest("POST", URITTS, bytes.NewBuffer(text))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEYTTS)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")

	// Print request *DEBUG
	res, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent Request:")
	fmt.Print(string(res))

	//Send POST request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	//Determine response
	if resp.StatusCode == http.StatusOK {
		t := map[string]interface{}{}
		if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
			fmt.Println("ERROR")
			log.Fatal(err)
			return nil, errors.New("Service")
		}
		return t["speech"], nil
	} else {
		fmt.Printf("\nStatus Code: \n %d\n", resp.StatusCode)
		log.Fatal(err)
		return nil, errors.New("Service")
	}
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
