package router

import (
	"log"
	"net/http"
)

func (rt Router) dataSourceRouteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		result, err := rt.dataSourceService.GetAllDataSource(r.Context())
		log.Printf("Fetched following: %+v\n", result)
		log.Printf("With error: %+v\n", err)

		rt.handleResult(result, err, w)
	case http.MethodPut:
		ds, err := rt.unmarshalBodyToDataSource(w, r)
		if err != nil {
			return
		}
		result, err := rt.dataSourceService.UpdateDataSource(r.Context(), ds)
		rt.handleResult(result, err, w)
	case http.MethodPost:
		ds, err := rt.unmarshalBodyToDataSource(w, r)
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
	id, err := rt.getPathParamHex(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Cannot parse the path param from(%s): %+v\n", r.URL.Path, err)
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
