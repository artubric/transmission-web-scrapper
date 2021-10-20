package main

import (
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/db"
	"transmission-web-scrapper/internal/server"
)

func main() {
	// fetch config
	config := config.GetConfig()

	// initiate connection to DB
	dbClient := db.Connect(config.DBConfig)

	// bring up rest APIs
	server := server.New(config.ServerConfig, dbClient)
	server.Run()
}
