package router

import (
	"log"
	"net/http"
)

func (rt Router) seasonRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		result, err := rt.seasonService.GetAllSeason(r.Context())
		rt.handleResult(result, err, w)
	case http.MethodPut:
		season, err := rt.unmarshalBodyToSeason(w, r)
		if err != nil {
			return
		}
		result, err := rt.seasonService.UpdateSeason(r.Context(), season)
		rt.handleResult(result, err, w)
	case http.MethodPost:
		season, err := rt.unmarshalBodyToSeason(w, r)
		if err != nil {
			return
		}
		result, err := rt.seasonService.CreateSeason(r.Context(), season)
		rt.handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rt Router) seasonRouteHandlerByID(w http.ResponseWriter, r *http.Request) {
	id, err := rt.getPathParam(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Cannot parse the path param from(%s): %+v\n", r.URL.Path, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		result, err := rt.seasonService.GetSeason(r.Context(), id)
		rt.handleResult(result, err, w)
	case http.MethodDelete:
		if err := rt.seasonService.DeleteSeason(r.Context(), id); err != nil {
			rt.writeErrorJSON(w, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
