package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

const (
	REGION = "uksouth"
	URI    = "https://" + REGION + ".tts.speech.microsoft.com/" +
		"cognitiveservices/v1"
	KEY = "d76745e51adf4408b1f29d7a4362dc39"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func TextToSpeech(text []byte) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", URI, bytes.NewBuffer(text))
	check(err)

	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", KEY)
	req.Header.Set("X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm")
	res, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sent Request:")
	fmt.Print(string(res))
	rsp, err2 := client.Do(req)
	check(err2)

	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusOK {
		body, err3 := ioutil.ReadAll(rsp.Body)
		check(err3)
		return body, nil
	} else {
		fmt.Printf("\nStatus Code: \n %d\n", rsp.StatusCode)
		return nil, errors.New("cannot convert text to speech")
	}
}

func main() {
	text, err := ioutil.ReadFile("text2.xml")
	check(err)
	speech, err2 := TextToSpeech(text)
	check(err2)
	err3 := ioutil.WriteFile("speech.wav", speech, 0644)
	check(err3)
}
