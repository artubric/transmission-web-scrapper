package main

import (
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/db"
	"transmission-web-scrapper/internal/router"
	"transmission-web-scrapper/internal/server"
	"transmission-web-scrapper/internal/service"
)

func main() {
	// fetch config
	config := config.Load()

	// initiate connection to DB
	dbRepositories := db.Connect(config.DBConfig)

	// services
	var telegramService *service.TelegramService
	if config.TelegramServiceConfig.Enabled {
		telegramService = service.NewTelegramService(config.TelegramServiceConfig)
	}
	sourceService := service.NewDataSourceService(dbRepositories.Source)
	seasonService := service.NewSeasonService(dbRepositories.Season)
	transmissionService := service.NewTransmissionService(config.TransmissionConfig)
	scraperService := service.NewScraperService(dbRepositories.Season, transmissionService, telegramService)

	// router
	router := router.New(sourceService, seasonService, scraperService)

	// bring up rest APIs
	server := server.New(config.ServerConfig, router)
	server.Run()
}
