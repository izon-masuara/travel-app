package repository

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationRepository interface {
	Save(ctx context.Context, db *mongo.Database, destination domain.DestinationCreate)
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
