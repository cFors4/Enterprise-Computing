package main

import (
	"encoding/json"
	"gorilla/mux"
	"net/http"
	"strconv"
)

var boxes map[string]string

// GET
func ListBoxes(w http.ResponseWriter, r *http.Request) {
	ids := []string{}
	for id, _ := range boxes {
		ids = append(ids, id)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ids)
}

// POST
func CreateBox(w http.ResponseWriter, r *http.Request) {
	// decode json
	id := strconv.Itoa(counter)
	counter = counter + 1
	t := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&t); err == nil {
		if cont, ok := t["contents"].(string); ok {
			w.Header().Set("Location", "/boxes/"+id)
			w.WriteHeader(http.StatusCreated) //201
			boxes[id] = cont
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

//GET ReadBox
func ReadBox(w http.ResponseWriter, r *http.Request) {

}

func main() {
	boxes = map[string]string{}
	var counter int
	counter = 0

	r := mux.NewRouter()
	r.HandleFunc("/boxes", CreateBox)

}

//PUT UpdateBox
//DELETE DeleteBox
