package main
import (
  "encoding/json"
  "errors"
  "net/http"
  "github.com/gorilla/mux"
)

const (
  URI = "http://api.what3words.com/v3/convert-to-3wa"
  KEY = "VZ2H2SZT"
)

func What3Words( w http.ResponseWriter, r *http.Request ) {
  t := map[ string ] interface{} {}
  if err := json.NewDecoder( r.Body ).Decode( &t ); err == nil {
    if coordinates, ok := t[ "coordinates" ].( string ); ok {
      if words, err := Service( coordinates ); err == nil {
        u := map[ string ] interface{} { "words" : words } 
        w.WriteHeader( http.StatusOK )
        json.NewEncoder( w ).Encode( u )
      } else {
        w.WriteHeader( http.StatusInternalServerError )
      }
    } else {
      w.WriteHeader( http.StatusBadRequest )
    }
  } else {
    w.WriteHeader( http.StatusBadRequest )
  }
}

func Service( coordinates string ) ( interface{}, error ) {
  client := &http.Client{}
  uri    := URI + "?key=" + KEY + "&coordinates=" + coordinates
  if req, err := http.NewRequest( "GET", uri, nil ); err == nil {
    if rsp, err := client.Do( req ); err == nil {
      if rsp.StatusCode == http.StatusOK {
        t := map[ string ] interface{} {} 
        if err := json.NewDecoder( rsp.Body ).Decode( &t ); err == nil {
          return t[ "words" ], nil
        }
      }
    }
  }
  return nil, errors.New( "Service" )
}

func main() {
  r := mux.NewRouter()
  // document 
  r.HandleFunc( "/what3words",  What3Words ).Methods( "POST" )
  http.ListenAndServe( ":3000", r )
}

