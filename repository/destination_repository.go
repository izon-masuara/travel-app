package repository

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationRepository interface {
	Save(ctx context.Context, db *mongo.Database, destination domain.DestinationCreate)
	FindAll(ctx context.Context, db *mongo.Database) []domain.Destination
}

type DestinationRepositoryImpl struct {
}

func NewDestinationRepository() DestinationRepository {
	return &DestinationRepositoryImpl{}
}

func (respository *DestinationRepositoryImpl) Save(ctx context.Context, db *mongo.Database, destination domain.DestinationCreate) {
	_, err := db.Collection("destination").InsertOne(ctx, destination)
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		helper.PanicIfError(err)
	}

}

func (repository *DestinationRepositoryImpl) FindAll(ctx context.Context, db *mongo.Database) []domain.Destination {
	cursor, err := db.Collection("destination").Find(ctx, bson.D{})
	helper.PanicIfError(err)
	defer cursor.Close(ctx)
	var result []domain.Destination
	for cursor.Next(ctx) {
		var row domain.Destination
		err := cursor.Decode(&row)
		helper.PanicIfError(err)
		result = append(result, row)
	}
	return result
}
