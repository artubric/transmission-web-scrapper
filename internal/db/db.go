package db

import (
	"context"
	"log"
	"time"

	"transmission-web-scrapper/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Connect(config config.DBConfig) *mongo.Client {
	log.Println("Connecting to db...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Failed to connect to DB with %+w", err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Failed to connect to DB with %+w", err)
	}

	log.Println("Done")

	return client
}
