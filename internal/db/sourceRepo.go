package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DataSourceRepository interface {
	Create(context.Context, DataSource) (DataSource, error)
	Delete(context.Context, primitive.ObjectID) error
	Get(context.Context, primitive.ObjectID) (DataSource, error)
	Update(context.Context, DataSource) (DataSource, error)
	GetAll(context.Context) ([]DataSource, error)
}

type DataSourceRepositoryImpl struct {
	col *mongo.Collection
}

func NewDataSourceRepository(col *mongo.Collection) DataSourceRepository {
	return &DataSourceRepositoryImpl{col: col}
}

func (dsr DataSourceRepositoryImpl) Get(ctx context.Context, id primitive.ObjectID) (DataSource, error) {
	var fetchedDataSource DataSource
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}

	result := dsr.col.FindOne(ctx, filter)
	if result.Err() != nil {
		return fetchedDataSource, result.Err()
	}
	if err := result.Decode(&fetchedDataSource); err != nil {
		return fetchedDataSource, err
	}

	return fetchedDataSource, nil
}

func (dsr DataSourceRepositoryImpl) GetAll(ctx context.Context) ([]DataSource, error) {
	var fetchedDataSource []DataSource
	var filter = bson.D{}

	result, err := dsr.col.Find(ctx, filter)
	if err != nil {
		return fetchedDataSource, err
	}
	if err := result.All(ctx, &fetchedDataSource); err != nil {
		return fetchedDataSource, err
	}

	return fetchedDataSource, nil
}

func (dsr DataSourceRepositoryImpl) Update(ctx context.Context, s DataSource) (DataSource, error) {

	if _, err := dsr.col.UpdateByID(ctx, s.ID, s); err != nil {
		return DataSource{}, err
	}

	return dsr.Get(ctx, s.ID)
}

func (dsr DataSourceRepositoryImpl) Create(ctx context.Context, s DataSource) (DataSource, error) {

	s.ID = primitive.NewObjectID()
	result, err := dsr.col.InsertOne(ctx, s)
	if err != nil {
		return DataSource{}, err
	}

	s.ID = result.InsertedID.(primitive.ObjectID)

	return s, nil
}

func (dsr DataSourceRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}

	if _, err := dsr.col.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
