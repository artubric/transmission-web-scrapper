package server

import (
	"fmt"
	"log"
	"net/http"
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/router"
)

type server struct {
	config config.ServerConfig
	router router.Router
}

func New(conf config.ServerConfig, router router.Router) server {
	return server{
		config: conf,
		router: router,
	}
}

func (s server) Run() {
	s.router.SetupRoutes(s.config.ApiBasePath)
	address := fmt.Sprintf(":%s", s.config.Port)
	log.Println("Starting server on: ", address)
	http.ListenAndServe(address, nil)
}
