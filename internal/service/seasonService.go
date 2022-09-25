package service

import (
	"context"
	"transmission-web-scrapper/internal/db"
	"transmission-web-scrapper/internal/dto"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SeasonService struct {
	repo db.SeasonRepository
}

func NewSeasonService(repo db.SeasonRepository) SeasonService {
	return SeasonService{
		repo: repo,
	}
}

func (s SeasonService) CreateSeason(ctx context.Context, season db.Season) (dto.Season, error) {
	result, err := s.repo.Create(ctx, season)
	if err != nil {
		return dto.Season{}, err
	}
	return dto.DbSeasonToDTO(result), nil
}

func (s SeasonService) UpdateSeason(ctx context.Context, season db.Season) (dto.Season, error) {
	result, err := s.repo.Update(ctx, season)
	if err != nil {
		return dto.Season{}, err
	}
	return dto.DbSeasonToDTO(result), nil
}

func (s SeasonService) DeleteSeason(ctx context.Context, id primitive.ObjectID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s SeasonService) GetSeason(ctx context.Context, id primitive.ObjectID) (dto.Season, error) {
	result, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.Season{}, err
	}
	return dto.DbSeasonToDTO(result), nil
}

func (s SeasonService) GetAllSeason(ctx context.Context, expandSource bool) ([]dto.Season, error) {
	dbSeasons, err := s.repo.GetAll(ctx, expandSource)
	if err != nil {
		return []dto.Season{}, err
	}
	return dto.DbSeasonsToDTO(dbSeasons), nil
}
