package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeasonRepository interface {
	Create(context.Context, Season) (Season, error)
	Delete(context.Context, primitive.ObjectID) error
	Get(context.Context, primitive.ObjectID) (Season, error)
	Update(context.Context, Season) (Season, error)
	GetAll(context.Context, bool) ([]Season, error)
}

type SeasonRepositoryImpl struct {
	col *mongo.Collection
}

func NewSeasonRepository(col *mongo.Collection) SeasonRepository {
	return &SeasonRepositoryImpl{col: col}
}

func (sr SeasonRepositoryImpl) Get(ctx context.Context, id primitive.ObjectID) (Season, error) {
	var fetchedSeason Season
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}

	result := sr.col.FindOne(ctx, filter)
	if result.Err() != nil {
		return fetchedSeason, result.Err()
	}
	if err := result.Decode(&fetchedSeason); err != nil {
		return fetchedSeason, err
	}

	return fetchedSeason, nil
}

func (sr SeasonRepositoryImpl) GetAll(ctx context.Context, expandSource bool) ([]Season, error) {
	var fetchedSeason []Season
	var result *mongo.Cursor
	var err error
	if expandSource {
		lookupStage := bson.D{primitive.E{
			Key: "$lookup",
			Value: bson.D{
				primitive.E{
					Key:   "from",
					Value: "Source",
				}, primitive.E{
					Key:   "localField",
					Value: "dataSourceId",
				}, primitive.E{
					Key:   "foreignField",
					Value: "_id",
				}, primitive.E{
					Key:   "as",
					Value: "dataSource",
				},
			},
		}}
		unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$dataSource"}}

		result, err = sr.col.Aggregate(ctx, mongo.Pipeline{
			lookupStage, unwindStage,
		})
		if err != nil {
			return fetchedSeason, err
		}
	} else {
		filter := bson.D{}
		result, err = sr.col.Find(ctx, filter)
		if err != nil {
			return fetchedSeason, err
		}
	}

	if err := result.All(ctx, &fetchedSeason); err != nil {
		return fetchedSeason, err
	}

	return fetchedSeason, nil
}

func (sr SeasonRepositoryImpl) Update(ctx context.Context, s Season) (Season, error) {
	var filter = bson.D{primitive.E{Key: "_id", Value: s.ID}}
	s.LastUpdated = primitive.NewDateTimeFromTime(time.Now())

	if _, err := sr.col.ReplaceOne(ctx, filter, s); err != nil {
		return Season{}, err
	}

	return sr.Get(ctx, s.ID)
}

func (sr SeasonRepositoryImpl) Create(ctx context.Context, s Season) (Season, error) {

	s.ID = primitive.NewObjectID()
	s.LastUpdated = primitive.NewDateTimeFromTime(time.Now())
	result, err := sr.col.InsertOne(ctx, s)
	if err != nil {
		return Season{}, err
	}

	s.ID = result.InsertedID.(primitive.ObjectID)

	return s, nil
}

func (sr SeasonRepositoryImpl) Delete(ctx context.Context, id primitive.ObjectID) error {
	var filter = bson.D{primitive.E{Key: "_id", Value: id}}

	if _, err := sr.col.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
