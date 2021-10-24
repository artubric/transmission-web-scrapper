package main

import (
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/db"
	"transmission-web-scrapper/internal/server"
)

func main() {
	// fetch config
	config := config.Load()

	// initiate connection to DB
	dbRepositories := db.Connect(config.DBConfig)

	// bring up rest APIs
	server := server.New(config.ServerConfig, dbRepositories)
	server.Run()
}
