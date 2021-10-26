package router

import (
	"io/ioutil"
	"net/http"
)

func (rt Router) scraperRouteHandler(w http.ResponseWriter, r *http.Request) {
	_, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPut:
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPost:
		w.WriteHeader(http.StatusNotImplemented)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
