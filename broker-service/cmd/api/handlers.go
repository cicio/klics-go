package main

import (
	"net/http"
)

// Create a handler function that as a receiver of type
// pointer to Config and takes two parameters of type
// http.ResponseWriter and pointer to http.Request
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	//specify a payload for testing purposes
	payload := jsonResponse{
		Error:   false,
		Message: "Send message to the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}
