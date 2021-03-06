package router

import (
	"fmt"
	"log"
	"net/http"
	"transmission-web-scrapper/internal/service"
)

type Router struct {
	seasonService     service.SeasonService
	dataSourceService service.DataSourceService
	scraperService    service.ScraperService
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func New(dataSourceService service.DataSourceService, seasonService service.SeasonService, scraperService service.ScraperService) Router {
	return Router{
		seasonService:     seasonService,
		dataSourceService: dataSourceService,
		scraperService:    scraperService,
	}
}

func (rt Router) SetupRoutes(apiBasePath string) {
	rt.newRoute(apiBasePath, "v1", "season", rt.seasonRouteHandler)
	rt.newRoute(apiBasePath, "v1", "season/", rt.seasonRouteHandlerByID)
	rt.newRoute(apiBasePath, "v1", "data-source", rt.dataSourceRouteHandler)
	rt.newRoute(apiBasePath, "v1", "data-source/", rt.dataSourceRouteHandlerByID)
	rt.newRoute(apiBasePath, "v1", "scraper", rt.scraperRouteHandler)
}

func (rt Router) newRoute(basePath string, apiVersion string, entityName string, hf handlerFunc) {
	urlPath := fmt.Sprintf("/%s/%s/%s",
		basePath,
		apiVersion,
		entityName,
	)

	log.Println("Registering listener for ", urlPath)
	http.Handle(urlPath, http.HandlerFunc(hf))
}
