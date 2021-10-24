package server

import (
	"fmt"
	"log"
	"net/http"
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/db"
)

type Server struct {
	config  config.ServerConfig
	dbRepos *db.DBRepositories
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func New(conf config.ServerConfig, dbRepos *db.DBRepositories) Server {
	return Server{
		config:  conf,
		dbRepos: dbRepos,
	}
}

func (s Server) Run() {
	s.setupRoutes()
	address := fmt.Sprintf(":%s", s.config.Port)
	log.Println("Starting server on: ", address)
	http.ListenAndServe(address, nil)
}

func (s *Server) setupRoutes() {
	newRoute(s.config.ApiBasePath, "v1", "season", s.seasonRouteHandler)
	newRoute(s.config.ApiBasePath, "v1", "data-source", s.dataSourceRouteHandler)
	newRoute(s.config.ApiBasePath, "v1", "scraper", s.scraperRouteHandler)
}

func newRoute(basePath string, apiVersion string, entityName string, hf handlerFunc) {
	urlPath := fmt.Sprintf("/%s/%s/%s",
		basePath,
		apiVersion,
		entityName,
	)

	log.Println("Registering listener for ", urlPath)
	http.Handle(urlPath, http.HandlerFunc(hf))
}
