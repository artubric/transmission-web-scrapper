package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Season struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name" bson:"name"`
	Season        int8               `json:"season" bson:"season"`
	StartDate     primitive.DateTime `json:"startDate" bson:"startDate"`
	EndDate       primitive.DateTime `json:"endDate" bson:"endDate"`
	TotalEpisodes uint16             `json:"totalEpisodes" bson:"totalEpisodes"`
	LastUpdated   primitive.DateTime `json:"lastUpdated" bson:"lastUpdated"`
	LastEpisode   uint16             `json:"lastEpisode" bson:"lastEpisode"`
	Quality       string             `json:"quality" bson:"quality"`
	DataSourceId  primitive.ObjectID `json:"dataSourceId" bson:"dataSourceId"`
	DataSource    DataSource         `json:"dataSource" bson:"dataSource"`
	IsArchived    bool               `json:"isArchived" bson:"isArchived"`
	DownloadDir   string             `json:"downloadDir" bson:"downloadDir"`
	ImdbId        string             `json:"imdbId" bson:"imdbId,omitempty"`
}

type DataSource struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Link       string             `json:"link" bson:"link"`
	SourceType string             `json:"sourceType" bson:"sourceType"`
	Parameters string             `json:"parameters" bson:"parameters"`
}
