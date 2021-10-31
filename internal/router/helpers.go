package router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"transmission-web-scrapper/internal/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (rt Router) writeErrorJSON(w http.ResponseWriter, err error) {
	errorJSON, err := json.Marshal(err)
	if err != nil {
		log.Printf("Error: %+v", err)
	}
	w.Write(errorJSON)
}

func (rt Router) unmarshallDataSource(body []byte) (db.DataSource, error) {
	var dataSource db.DataSource
	if err := json.Unmarshal(body, &dataSource); err != nil {
		return db.DataSource{}, err
	}
	return dataSource, nil
}

func (rt Router) unmarshallSeason(body []byte) (db.Season, error) {
	var season db.Season
	if err := json.Unmarshal(body, &season); err != nil {
		return db.Season{}, err
	}
	return season, nil
}

func (rt Router) unmarshalBodyToDataSource(w http.ResponseWriter, r *http.Request) (db.DataSource, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rt.writeErrorJSON(w, err)
		log.Printf("Cannot parse the body: %+v\n", err)
		return db.DataSource{}, err
	}
	ds, err := rt.unmarshallDataSource(body)
	if err != nil {
		log.Printf("Cannot unmarshal the body: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		rt.writeErrorJSON(w, err)
		return db.DataSource{}, err
	}
	return ds, nil
}

func (rt Router) unmarshalBodyToSeason(w http.ResponseWriter, r *http.Request) (db.Season, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rt.writeErrorJSON(w, err)
		log.Printf("Cannot parse the body: %+v\n", err)
		return db.Season{}, err
	}
	ds, err := rt.unmarshallSeason(body)
	if err != nil {
		log.Printf("Cannot unmarshal the body: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		rt.writeErrorJSON(w, err)
		return db.Season{}, err
	}
	return ds, nil
}

func (rt Router) marshalObject(i interface{}) ([]byte, error) {
	resultJSON, err := json.Marshal(&i)
	if err != nil {
		return nil, err
	}
	return resultJSON, nil
}

func (rt Router) handleResult(i interface{}, err error, w http.ResponseWriter) {
	if err != nil {
		log.Printf("Failed with: %+v\n", err)
		rt.writeErrorJSON(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	jsonResponse, err := rt.marshalObject(i)
	if err != nil {
		rt.writeErrorJSON(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write(jsonResponse)
	w.WriteHeader(http.StatusOK)
}

func (rt Router) getPathParam(path string) (primitive.ObjectID, error) {
	pathSlice := strings.Split(path, "/")
	inputId := pathSlice[len(pathSlice)-1]
	return primitive.ObjectIDFromHex(inputId)
}
