package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	//a variable that has email an password tags
	//that we are going to receive in a JSON.
	//The  received data will be unmask in the requestPayload
	var requestPayload struct {
		Email    string `json:"email`
		Password string `json:"password`
	}
	//check for errors when I call the readJSON function from helprs.go
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)

	}

	//validate the user against the database
	// using user.passwordMatches
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentilas"), http.StatusBadRequest)

	}
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentilas"), http.StatusBadRequest)
		return
	}
	// provides the payload of type jsonResponse
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}
