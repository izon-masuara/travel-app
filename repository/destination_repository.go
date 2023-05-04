package repository

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationRepository interface {
	Save(ctx context.Context, db *mongo.Database, destination domain.DestinationCreate)
	FindAll(ctx context.Context, db *mongo.Database) []domain.Destination
	Update(ctx context.Context, db *mongo.Database, destination domain.DestinationUpdate, requestId string)
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

func (repository *DestinationRepositoryImpl) Update(ctx context.Context, db *mongo.Database, destination domain.DestinationUpdate, requestId string) {
	id, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		panic(exception.NewNotFoundError("Data not found"))
	}
	var found domain.Destination
	err = db.Collection("destination").FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&found)
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		panic(exception.NewNotFoundError("Data not found"))
	}
	_, err = db.Collection("destination").UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"title":      destination.Title,
			"date":       destination.Date,
			"image_file": destination.ImageFile,
			"location": bson.M{
				"long": destination.Location.Long,
				"lat":  destination.Location.Lat,
			},
			"text": destination.Text,
		},
	})
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		helper.PanicIfError(err)
	}
	helper.RemoveFile(found.ImageFile)
}
