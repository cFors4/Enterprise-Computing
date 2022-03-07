package main
import (
  "encoding/json"
  "net/http"
  "strconv"
  "github.com/gorilla/mux"
)

var boxes map[ string ] string

var counter int

func CreateBox( w http.ResponseWriter, r *http.Request ) {
  id := strconv.Itoa( counter ); counter = counter + 1
  t  := map[ string ] string {}
  if err := json.NewDecoder( r.Body ).Decode( &t ); err == nil {
    if cont, ok := t[ "contents" ]; ok {
      w.Header().Set( "Location", "/boxes/" + id )
      w.WriteHeader( http.StatusCreated )
      boxes[ id ] = cont
    } else {
      w.WriteHeader( http.StatusBadRequest )
    }
  } else {
    w.WriteHeader( http.StatusBadRequest )
  }
}

func ListBoxes( w http.ResponseWriter, r *http.Request ) {
  ids := []string{}
  for id, _ := range boxes {
    ids = append( ids, id )
  }
  w.WriteHeader( http.StatusOK )
  json.NewEncoder( w ).Encode( ids )
}

func ReadBox( w http.ResponseWriter, r *http.Request ) {
  vars := mux.Vars( r )
  id   := vars[ "id" ]
  if contents, ok := boxes[ id ]; ok {
    u := map[ string ] string { "contents" : contents }
    w.WriteHeader( http.StatusOK )
    json.NewEncoder( w ).Encode( u )
  } else {
    w.WriteHeader( http.StatusNotFound )
  }
}

func main() {
  boxes   = make( map [ string ] string )
  counter = 0
  r := mux.NewRouter()
  // collection
  r.HandleFunc( "/boxes",      CreateBox ).Methods( "POST"   )
  r.HandleFunc( "/boxes",      ListBoxes ).Methods( "GET"    )
  // document
  r.HandleFunc( "/boxes/{id}", ReadBox   ).Methods( "GET"    )

  http.ListenAndServe( ":3000", r )
}

