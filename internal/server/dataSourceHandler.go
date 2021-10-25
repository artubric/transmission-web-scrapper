package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"transmission-web-scrapper/internal/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Server) dataSourceRouteHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		result, err := s.dataSourceService.GetAllDataSource(r.Context())
		log.Printf("Fetched following: %+v\n", result)
		log.Printf("With error: %+v\n", err)

		handleResult(result, err, w)
	case http.MethodPut:
		ds, err := unmarshalBody(w, r)
		if err != nil {
			return
		}
		result, err := s.dataSourceService.UpdateDataSource(r.Context(), ds)
		handleResult(result, err, w)
	case http.MethodPost:
		ds, err := unmarshalBody(w, r)
		if err != nil {
			return
		}
		result, err := s.dataSourceService.CreateDataSource(r.Context(), ds)
		handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s Server) dataSourceRouteHandlerByID(w http.ResponseWriter, r *http.Request) {
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
		result, err := s.dataSourceService.GetDataSource(r.Context(), id)
		handleResult(result, err, w)
	case http.MethodDelete:
		if err := s.dataSourceService.DeleteDataSource(r.Context(), id); err != nil {
			writeErrorJSON(w, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func writeErrorJSON(w http.ResponseWriter, err error) {
	errorJSON, err := json.Marshal(err)
	if err != nil {
		log.Printf("Error: %+v", err)
	}
	log.Printf("Error json: %+v", errorJSON)

	w.Write(errorJSON)
}

func unmarshallDataSource(body []byte) (db.DataSource, error) {
	var dataSource db.DataSource
	if err := json.Unmarshal(body, &dataSource); err != nil {
		return dataSource, err
	}
	return dataSource, nil
}

func unmarshalBody(w http.ResponseWriter, r *http.Request) (db.DataSource, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("Cannot parse the body: %+v\n", err)
		return db.DataSource{}, err
	}
	ds, err := unmarshallDataSource(body)
	if err != nil {
		log.Printf("Cannot unmarshal the body: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		writeErrorJSON(w, err)
		return db.DataSource{}, err
	}
	return ds, nil
}

func marshalObject(i interface{}) ([]byte, error) {
	resultJSON, err := json.Marshal(&i)
	if err != nil {
		return nil, err
	}
	return resultJSON, nil
}

func handleResult(i interface{}, err error, w http.ResponseWriter) {
	if err != nil {
		writeErrorJSON(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonResponse, err := marshalObject(i)
	if err != nil {
		writeErrorJSON(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}
