package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// Create a handler function that as a receiver of type
// pointer to Config and takes two parameters of type
// http.ResponseWriter and pointer to http.Request
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	//specify a payload for testing purposes
	payload := jsonResponse{
		Error:   false,
		Message: "Send message to the broker",
	}

	//write this data out. Manually for now
	out, _ := json.MarshalIndent(payload, "", "\t")    //MarshalIndent beaufifies it
	w.Header().Set("Content-Type", "application/json") //set the Header
	w.WriteHeader(http.StatusAccepted)                 // add the Header
	w.Write(out)                                       //send out the payload

}
