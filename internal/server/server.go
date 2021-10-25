package server

import (
	"fmt"
	"log"
	"net/http"
	"transmission-web-scrapper/config"
	"transmission-web-scrapper/internal/service"
)

type Server struct {
	config            config.ServerConfig
	seasonService     *service.SeasonService
	dataSourceService *service.DataSourceService
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func New(conf config.ServerConfig, seasonService *service.SeasonService, dataSourceService *service.DataSourceService) Server {
	return Server{
		config:            conf,
		seasonService:     seasonService,
		dataSourceService: dataSourceService,
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
	newRoute(s.config.ApiBasePath, "v1", "data-source/", s.dataSourceRouteHandlerByID)
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
