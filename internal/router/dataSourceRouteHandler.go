package router

import (
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (rt Router) dataSourceRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		result, err := rt.dataSourceService.GetAllDataSource(r.Context())
		log.Printf("Fetched following: %+v\n", result)
		log.Printf("With error: %+v\n", err)

		rt.handleResult(result, err, w)
	case http.MethodPut:
		ds, err := rt.unmarshalBody(w, r)
		if err != nil {
			return
		}
		result, err := rt.dataSourceService.UpdateDataSource(r.Context(), ds)
		rt.handleResult(result, err, w)
	case http.MethodPost:
		ds, err := rt.unmarshalBody(w, r)
		if err != nil {
			return
		}
		result, err := rt.dataSourceService.CreateDataSource(r.Context(), ds)
		rt.handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rt Router) dataSourceRouteHandlerByID(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	inputId := path[len(path)-1]
	id, err := primitive.ObjectIDFromHex(inputId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Cannot parse the path param(%+v): %+v\n", inputId, err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		result, err := rt.dataSourceService.GetDataSource(r.Context(), id)
		rt.handleResult(result, err, w)
	case http.MethodDelete:
		if err := rt.dataSourceService.DeleteDataSource(r.Context(), id); err != nil {
			rt.writeErrorJSON(w, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
