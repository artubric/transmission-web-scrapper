package dto

import (
	"fmt"
	"time"
	"transmission-web-scrapper/internal/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Season struct {
	ID            string `json:"id"`
	Name          string `json:"name" `
	Season        int8   `json:"season"`
	StartDate     string `json:"startDate"`
	EndDate       string `json:"endDate"`
	TotalEpisodes uint16 `json:"totalEpisodes"`
	LastUpdated   string `json:"lastUpdated"`
	LastEpisode   uint16 `json:"lastEpisode"`
	Quality       string `json:"quality"`
	DataSourceId  string `json:"dataSourceId"`
	IsArchived    bool   `json:"isArchived"`
	DownloadDir   string `json:"downloadDir"`
}

func DTOSeasonToDB(s Season) (db.Season, error) {

	var dbId primitive.ObjectID
	var err error
	if s.ID == "add" {
		dbId = primitive.NilObjectID
	} else {
		dbId, err = primitive.ObjectIDFromHex(s.ID)
		if err != nil {
			return db.Season{}, fmt.Errorf("Failed to parse ID to primitive, %w", err)
		}
	}

	dbDataSourceId, err := primitive.ObjectIDFromHex(s.DataSourceId)
	if err != nil {
		return db.Season{}, fmt.Errorf("Failed to parse DataSourceId to primitive, %w", err)
	}

	startTime, err := time.Parse("2006-01-02T15:04:05", s.StartDate)
	if err != nil {
		return db.Season{}, fmt.Errorf("Failed to parse StartDate to time, %w", err)
	}

	endTime, err := time.Parse("2006-01-02T15:04:05", s.EndDate)
	if err != nil {
		return db.Season{}, fmt.Errorf("Failed to parse EndDate to time, %w", err)
	}
	return db.Season{
		ID:            dbId,
		Name:          s.Name,
		Season:        s.Season,
		StartDate:     primitive.NewDateTimeFromTime(startTime),
		EndDate:       primitive.NewDateTimeFromTime(endTime),
		TotalEpisodes: s.TotalEpisodes,
		LastUpdated:   primitive.NewDateTimeFromTime(time.Now()),
		LastEpisode:   s.LastEpisode,
		Quality:       s.Quality,
		DataSourceId:  dbDataSourceId,
		IsArchived:    s.IsArchived,
		DownloadDir:   s.DownloadDir,
	}, nil
}

func DbSeasonToDTO(s db.Season) Season {
	return Season{
		ID:            s.ID.Hex(),
		Name:          s.Name,
		Season:        s.Season,
		StartDate:     s.StartDate.Time().Format("2006-01-02T15:04:05"),
		EndDate:       s.EndDate.Time().Format("2006-01-02T15:04:05"),
		TotalEpisodes: s.TotalEpisodes,
		LastUpdated:   s.LastUpdated.Time().Format("2006-01-02T15:04:05"),
		LastEpisode:   s.LastEpisode,
		Quality:       s.Quality,
		DataSourceId:  s.DataSourceId.Hex(),
		IsArchived:    s.IsArchived,
		DownloadDir:   s.DownloadDir,
	}
}

func DbSeasonsToDTO(dbSeasons []db.Season) []Season {
	dtoSeasons := []Season{}
	for _, dbSeason := range dbSeasons {
		dtoSeasons = append(dtoSeasons, DbSeasonToDTO(dbSeason))
	}

	return dtoSeasons
}
