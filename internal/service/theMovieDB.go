package service

import (
	"fmt"
	"transmission-web-scrapper/config"

	tmdb "github.com/cyruzin/golang-tmdb"
)

type TmdbAPI interface {
	SearchTVShow(query string) (*tmdb.SearchTVShowsResults, error)
	GetTVShowDetails(tmdbShowId int) (*tmdb.TVDetails, error)
}

type tmdbService struct {
	client *tmdb.Client
}


func NewTMDBService(tmdbConfig config.TmdbAPIConfig) TmdbAPI {
	tmdbClient, err := tmdb.Init(tmdbConfig.ApiKey)
	if err != nil {
		fmt.Println(err)
	}

	return tmdbService{
		client: tmdbClient,
	}
}

func (s tmdbService) SearchTVShow(query string) (*tmdb.SearchTVShowsResults, error) {
	result, err := s.client.GetSearchTVShow(query, nil)
	if err != nil {
		return nil, err
	}
	return result.SearchTVShowsResults, nil
}

func (s tmdbService) GetTVShowDetails(tmdbShowId int) (*tmdb.TVDetails, error){
	return s.client.GetTVDetails(tmdbShowId, nil)
}


