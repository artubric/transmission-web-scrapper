package service

import (
	"context"
	"transmission-web-scrapper/internal/db"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DataSourceService struct {
	repo db.DataSourceRepository
}

func NewDataSourceService(repo db.DataSourceRepository) DataSourceService {
	return DataSourceService{
		repo: repo,
	}
}

func (s DataSourceService) CreateDataSource(ctx context.Context, ds db.DataSource) (db.DataSource, error) {
	result, err := s.repo.Create(ctx, ds)
	if err != nil {
		return db.DataSource{}, err
	}
	return result, nil
}

func (s DataSourceService) UpdateDataSource(ctx context.Context, ds db.DataSource) (db.DataSource, error) {
	result, err := s.repo.Update(ctx, ds)
	if err != nil {
		return db.DataSource{}, err
	}
	return result, nil
}

func (s DataSourceService) DeleteDataSource(ctx context.Context, id primitive.ObjectID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s DataSourceService) GetDataSource(ctx context.Context, id primitive.ObjectID) (db.DataSource, error) {
	result, err := s.repo.Get(ctx, id)
	if err != nil {
		return db.DataSource{}, err
	}
	return result, nil
}

func (s DataSourceService) GetAllDataSource(ctx context.Context) ([]db.DataSource, error) {
	result, err := s.repo.GetAll(ctx)
	if err != nil {
		return []db.DataSource{}, err
	}
	return result, nil
}
