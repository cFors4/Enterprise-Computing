package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	URIALPHA = "https://api.wolframalpha.com/v1/"
	APPID    = "JXWPV6-WW86KHYYLR"
)

type request_struct struct {
	Text string `json:"text"`
}

func Alpha(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/alpha" {
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
		if output, err := ServiceAlpha(text); err == nil {
			u := map[string]interface{}{"text": output}
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

func ServiceAlpha(input string) (string, error) {
	appID := "JXWPV6-WW86KHYYLR"
	input = strings.ReplaceAll(input, " ", "+")
	u := URIALPHA + "result?appid=" + appID + "&i=" + input + "%3F&timeout=5&units=metric"
	fmt.Printf("%s\n", u)

	resp, err := http.Get(u)
	if err != nil {
		return "", errors.New("Cannot GET response from wolfram alpha")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Cannot read response from wolfram alpha")
	}

	return string(body), nil
}

func main() {
	port := 3001
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/alpha", Alpha)
	fmt.Printf("Listening on port: %d\n", port)
	if err := http.ListenAndServe(portStr, nil); err != nil {
		log.Fatal(err)
	}
}
