package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// type RequestPayload will define the struc for
// json request that the handle will process
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

// AuthPayload defines the struct with the parameters for an authentication
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// logPayload defines the structure for a log entry
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
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

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	//create a variable of type RequestPayload to
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// Depending on the content received in the JSON payload
	// we want to take a different action. make use of the switch statement.
	// We want to switch in the requestPayload.Action
	switch requestPayload.Action {
	case "auth":
		//use the authenticate function to factor out code for simplicity
		app.Authenticate(w, requestPayload.Auth)

	case "log":
		app.logItem(w, requestPayload.Log)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) logItem(w http.ResponseWriter, entry LogPayload) {
	//create some json we'll send to the logger-service. do not use MarshallIndent in production code
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	//call the service
	request, err := http.NewRequest(
		"POST",
		logServiceURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service
	request, err := http.NewRequest(
		"POST",
		"http://authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling the service"))
		return
	}

	// create a variable we'll read response body into
	var jsonFromService jsonResponse

	// decode the json from auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	//create the payload to send back to the requester
	// since no errors were detected we will say Error = false and
	// the approriate message = Authenticated and the data to be sent is the
	// jsonFromservice.Data
	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated"
	payload.Data = jsonFromService.Data

	// Write the payload value to json as wella sthe http status Accepted
	app.writeJSON(w, http.StatusAccepted, payload)

}
