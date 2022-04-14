package dto

import "transmission-web-scrapper/internal/db"

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
