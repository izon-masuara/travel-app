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
	Destroy(ctx context.Context, db *mongo.Database, requestId string)
	FindOne(ctx context.Context, db *mongo.Database, requestId string) domain.Destination
}

type DestinationRepositoryImpl struct {
}

func NewDestinationRepository() DestinationRepository {
	return &DestinationRepositoryImpl{}
}

func (respository *DestinationRepositoryImpl) Save(ctx context.Context, db *mongo.Database, destination domain.DestinationCreate) {
	jsonAuth := helper.InterfaceToJsonAuth(ctx)
	if jsonAuth.Role != "operator" {
		panic(exception.NewAuthError("UNAUTORIZED"))
	}
	_, err := db.Collection(jsonAuth.Name).InsertOne(ctx, destination)
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		helper.PanicIfError(err)
	}

}

func (repository *DestinationRepositoryImpl) FindAll(ctx context.Context, db *mongo.Database) []domain.Destination {
	jsonAuth := helper.InterfaceToJsonAuth(ctx)
	if jsonAuth.Role != "operator" {
		panic(exception.NewAuthError("UNAUTORIZED"))
	}
	cursor, err := db.Collection(jsonAuth.Name).Find(ctx, bson.D{})
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
	jsonAuth := helper.InterfaceToJsonAuth(ctx)
	if jsonAuth.Role != "operator" {
		panic(exception.NewAuthError("UNAUTORIZED"))
	}
	id, err := primitive.ObjectIDFromHex(requestId)
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		panic(exception.NewNotFoundError("Data not found"))
	}
	var found domain.Destination
	err = db.Collection(jsonAuth.Name).FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&found)
	if err != nil {
		helper.RemoveFile(destination.ImageFile)
		panic(exception.NewNotFoundError("Data not found"))
	}
	_, err = db.Collection(jsonAuth.Name).UpdateOne(ctx, bson.M{"_id": id}, bson.M{
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

func (repository *DestinationRepositoryImpl) Destroy(ctx context.Context, db *mongo.Database, requestId string) {
	jsonAuth := helper.InterfaceToJsonAuth(ctx)
	if jsonAuth.Role != "operator" {
		panic(exception.NewAuthError("UNAUTORIZED"))
	}
	id, err := primitive.ObjectIDFromHex(requestId)
	helper.PanicIfError(err)
	var found domain.Destination
	err = db.Collection(jsonAuth.Name).FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&found)
	if err != nil {
		panic(exception.NewNotFoundError("Data not found"))
	}
	_, err = db.Collection(jsonAuth.Name).DeleteOne(ctx, bson.M{"_id": id})
	helper.PanicIfError(err)
	helper.RemoveFile(found.ImageFile)
}

func (repository *DestinationRepositoryImpl) FindOne(ctx context.Context, db *mongo.Database, requestId string) domain.Destination {
	jsonAuth := helper.InterfaceToJsonAuth(ctx)
	if jsonAuth.Role != "operator" {
		panic(exception.NewAuthError("UNAUTORIZED"))
	}
	id, err := primitive.ObjectIDFromHex(requestId)
	helper.PanicIfError(err)
	var found domain.Destination
	err = db.Collection(jsonAuth.Name).FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&found)
	if err != nil {
		panic(exception.NewNotFoundError("Data not found"))
	}
	return found
}
