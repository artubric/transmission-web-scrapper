package router

import (
	"net/http"
)

func (rt Router) scraperRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if err := rt.scraperService.Start(r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			rt.writeErrorJSON(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
