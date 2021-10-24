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

const (
	dBName               = "TransmissionWebScrapper" // TODO: config?
	seasonCollectionName = "Season"
	sourceCollectionName = "Source"
	timeout              = 10 * time.Second // TODO: config?
)

type DBRepositories struct {
	Season SeasonRepository
	Source DataSourceRepository
}

func Connect(config config.DBConfig) *DBRepositories {
	log.Println("Connecting to db...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.URI))
	if err != nil {
		log.Fatal("Failed to connect to DB with %+w", err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Failed to connect to DB with %+w", err)
	}

	dbRepositories := setupRepositories(client)

	return &dbRepositories

}

func setupRepositories(client *mongo.Client) DBRepositories {
	seasonCollection := client.Database(dBName).Collection(seasonCollectionName)
	sourceCollection := client.Database(dBName).Collection(sourceCollectionName)
	seasonRepository := NewSeasonRepository(seasonCollection)
	sourceRepository := NewDataSourceRepository(sourceCollection)
	return DBRepositories{
		Season: seasonRepository,
		Source: sourceRepository,
	}
}
