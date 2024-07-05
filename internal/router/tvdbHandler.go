package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"transmission-web-scrapper/internal/dto"

	tmdb "github.com/cyruzin/golang-tmdb"
)

func (rt Router) tmdbTVShowSearchHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		searchQuery := r.URL.Query().Get("searchQuery")
		if len(searchQuery) == 0 {
			err := fmt.Errorf("empty 'searchQuery' query parameter")
			rt.handleResult(tmdb.SearchTVShowsResults{}, err, w)
		}
		result, err := rt.tmdbAPIService.SearchTVShow(searchQuery)
		rt.handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rt Router) tmdbTVShowHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		result, err := rt.getTVShowDetails(r.URL.Path)
		rt.handleResult(result, err, w)

	case http.MethodPost:
		tvShowDetails, err := rt.getTVShowDetails(r.URL.Path)
		if err != nil {
			rt.handleResult(nil, err, w)
		}

		dataSourceId := r.URL.Query().Get("dataSourceId")
		if len(dataSourceId) == 0 {
			err := fmt.Errorf("dataSourceId query param is missing")
			rt.handleResult(nil, err, w)
			return
		}

		seasonNumberString := r.URL.Query().Get("seasonNumber")
		if len(seasonNumberString) == 0 {
			err := fmt.Errorf("seasonNumber query param is missing")
			rt.handleResult(nil, err, w)
			return
		}

		seasonNumber, err := strconv.Atoi(seasonNumberString)
		if err != nil {
			rt.handleResult(nil, err, w)
			return
		}

		var SeasonNumber int8
		var EpisodeCount uint16
		found := false

		for _, tvshowSeason := range tvShowDetails.Seasons {
			if tvshowSeason.SeasonNumber == seasonNumber {
				SeasonNumber = int8(tvshowSeason.SeasonNumber)
				EpisodeCount = uint16(tvshowSeason.EpisodeCount)
				found = true
				break
			}
		}

		if !found {
			err := fmt.Errorf("did not found requested season")
			rt.handleResult(nil, err, w)
			return
		}

		dtoSeason := dto.Season{
			ID:        "add",
			Name:      tvShowDetails.Name,
			Season:    SeasonNumber,
			StartDate: time.Now().Format("2006-01-02T15:04:05"),
			EndDate: time.Now().
				Add(2200 * time.Hour). // around 3 months
				Format("2006-01-02T15:04:05"),
			TotalEpisodes: EpisodeCount,
			Quality:       "1080p",
			DataSourceId:  dataSourceId,
			IsArchived:    false,
			DownloadDir:   fmt.Sprintf("/plex/tv_shows/%s/s%02d", strings.ToLower(tvShowDetails.Name), SeasonNumber),
			LastEpisode:   0,
			//ImdbId:        tvShowDetails.IMDbID, missing from response
		}

		dbSeason, err := dto.DTOSeasonToDB(dtoSeason)
		if err != nil {
			rt.handleResult(nil, err, w)
			return
		}

		result, err := rt.seasonService.CreateSeason(r.Context(), dbSeason)
		rt.handleResult(result, err, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (rt Router) getTVShowDetails(urlPath string) (*tmdb.TVDetails, error) {
	showIdString := rt.getPathParamString(urlPath)

	showId, err := strconv.Atoi(showIdString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse showIdString %w", err)
	}

	return rt.tmdbAPIService.GetTVShowDetails(showId)
}
