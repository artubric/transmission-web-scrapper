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
	tmdbAPIService    service.TmdbAPI
}

type handlerFunc func(w http.ResponseWriter, r *http.Request)

func New(dataSourceService service.DataSourceService, seasonService service.SeasonService, scraperService service.ScraperService, tmdbAPIService service.TmdbAPI) Router {
	return Router{
		seasonService:     seasonService,
		dataSourceService: dataSourceService,
		scraperService:    scraperService,
		tmdbAPIService:    tmdbAPIService,
	}
}

func (rt Router) SetupRoutes(apiBasePath string) {
	rt.newRoute(apiBasePath, "v1", "season", rt.seasonRouteHandler)
	rt.newRoute(apiBasePath, "v1", "season/", rt.seasonRouteHandlerByID)
	rt.newRoute(apiBasePath, "v1", "data-source", rt.dataSourceRouteHandler)
	rt.newRoute(apiBasePath, "v1", "data-source/", rt.dataSourceRouteHandlerByID)
	rt.newRoute(apiBasePath, "v1", "scraper", rt.scraperRouteHandler)
	rt.newRoute(apiBasePath, "v1", "tv-show/search", rt.tmdbTVShowSearchHandler)
	rt.newRoute(apiBasePath, "v1", "tv-show/", rt.tmdbTVShowHandler)

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
