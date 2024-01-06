package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"transmission-web-scrapper/internal/db"
	"transmission-web-scrapper/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (rt Router) writeErrorJSON(w http.ResponseWriter, err error) {
	log.Printf("Failed with: %+v", err)
	errorJSON, err := json.Marshal(err)
	if err != nil {
		log.Printf("Error while marshalling error: %+v", err)
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
	var dtoSeason dto.Season
	if err := json.Unmarshal(body, &dtoSeason); err != nil {
		return db.Season{}, err
	}
	dbSeason, err := dto.DTOSeasonToDB(dtoSeason)
	if err != nil {
		return db.Season{}, err
	}
	return dbSeason, nil
}

func (rt Router) unmarshalBodyToDataSource(w http.ResponseWriter, r *http.Request) (db.DataSource, error) {
	body, err := io.ReadAll(r.Body)
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
	body, err := io.ReadAll(r.Body)
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
		w.WriteHeader(http.StatusBadRequest)
		rt.writeErrorJSON(w, err)
		return
	}
	jsonResponse, err := rt.marshalObject(i)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		rt.writeErrorJSON(w, err)
		return
	}
	w.Write(jsonResponse)
}

func (rt Router) getPathParamHex(path string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(rt.getPathParamString(path))
}

func (rt Router) getPathParamString(path string) (string) {
	pathSlice := strings.Split(path, "/")
	inputId := pathSlice[len(pathSlice)-1]
	return inputId
}
