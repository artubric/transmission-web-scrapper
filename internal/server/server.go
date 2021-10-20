package server

import (
	"context"
	"log"
	"transmission-web-scrapper/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Server struct {
	config   config.ServerConfig
	dbClient *mongo.Client
}

func New(conf config.ServerConfig, dbClient *mongo.Client) Server {
	return Server{
		config:   conf,
		dbClient: dbClient,
	}
}

func (s Server) Run() {
	log.Print("Oh yea, running the server!")
	if err := s.dbClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal("Failed to ping the db")
	}
	log.Print("Succesfully pinged the DB")
}
