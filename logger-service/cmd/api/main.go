package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cicio/klics-go/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {

	//connect to Mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)

	}
	//assign client the value of mongoClient variable
	client = mongoClient

	// create a context in order to disconnect needed by mongo
	ctx, cancel := context.WithTimeout((context.Background()), 15*time.Second)
	defer cancel()

	//close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start web server as a go routine
	log.Println("Starting webserver on port", webPort)
	go app.serve()

}

// function to start a webserver
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

func connectToMongo() (*mongo.Client, error) {
	//create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	//connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting", err)
		return nil, err
	}

	log.Println("Connected to Mongo")
	return c, nil
}
