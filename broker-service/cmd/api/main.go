package main

import (
	"fmt"
	"log"
	"net/http"
)

// we will be using docker.
// And docker will listen on port 80 for any container
const webport = "8081"

// declare a type config of type struct that will be receiver
// for the application
type Config struct{}

func main() {
	// create a variable `app` of type config
	app := Config{}

	//create log to print
	log.Printf("Starting broker service on port %s\n", webport)

	//define an http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webport),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
