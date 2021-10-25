package service

import (
	"transmission-web-scrapper/internal/db"
)

type SeasonService struct {
	repo db.SeasonRepository
}

func NewSeasonService(repo db.SeasonRepository) SeasonService {
	return SeasonService{
		repo: repo,
	}
}
