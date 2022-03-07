package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	URIALPHA = "https://api.wolframalpha.com/v1/"
	APPID    = "JXWPV6-WW86KHYYLR"
)

func Alpha(w http.ResponseWriter, r *http.Request) {
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if question, ok := t["text"].(string); ok {
			if text, err := ServiceAlpha(question); err == nil {
				u := map[string]interface{}{"text": text}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(u)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ServiceAlpha(question string) (interface{}, error) {
	client := &http.Client{}
	uri := URIALPHA + "result?i=" + question + "%3F&appid=" + APPID
	if req, err := http.NewRequest("GET", uri, nil); err == nil {
		if rsp, err := client.Do(req); err == nil {
			if rsp.StatusCode == http.StatusOK {
				t := map[string]interface{}{}
				if err := json.NewDecoder(rsp.Body).Decode(&t); err == nil {
					return t["text"], nil
				}
			}
		}
	}
	return nil, errors.New("Service")
}

func main() {
	r := mux.NewRouter()
	// document
	r.HandleFunc("/alpha", Alpha).Methods("POST")
	http.ListenAndServe(":3001", r)
}
