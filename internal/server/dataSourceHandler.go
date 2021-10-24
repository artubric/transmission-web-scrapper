package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"transmission-web-scrapper/internal/db"
)

func (s Server) dataSourceRouteHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
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
		var dataSource db.DataSource
		var err = json.Unmarshal(body, &dataSource)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorJSON, _ := json.Marshal(err)
			w.Write(errorJSON)
			log.Println(err)
			return
		}

		created, err := s.dbRepos.Source.Create(r.Context(), dataSource)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorJSON, _ := json.Marshal(err)
			w.Write(errorJSON)
			return
		}

		createdJSON, err := json.Marshal(created)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorJSON, _ := json.Marshal(err)
			w.Write(errorJSON)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(createdJSON)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
