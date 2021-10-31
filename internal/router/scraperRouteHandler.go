package router

import (
	"net/http"
)

func (rt Router) scraperRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		rt.scraperService.Start(r.Context())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
