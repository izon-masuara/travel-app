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

type OperatorRepository interface {
	Save(ctx context.Context, db *mongo.Database, operator domain.OperatorCreate)
	FindAll(ctx context.Context, db *mongo.Database) []domain.OperatorSchema
	ResetPasswordById(ctx context.Context, db *mongo.Database, operatorId string)
	Destroy(ctx context.Context, db *mongo.Database, operatorId string)
}

type OperatorRepositoryImpl struct {
}

func NewOperatorRepository() OperatorRepository {
	return &OperatorRepositoryImpl{}
}

func (repository *OperatorRepositoryImpl) Save(ctx context.Context, db *mongo.Database, operator domain.OperatorCreate) {
	hash, err := helper.HashPassword(operator.Password)
	helper.PanicIfError(err)
	operator.Password = hash
	_, err = db.Collection("student").InsertOne(ctx, operator)
	helper.PanicIfError(err)
}

func (repository *OperatorRepositoryImpl) FindAll(ctx context.Context, db *mongo.Database) []domain.OperatorSchema {
	cursor, err := db.Collection("student").Find(ctx, bson.D{})
	helper.PanicIfError(err)
	defer cursor.Close(ctx)
	var result []domain.OperatorSchema
	for cursor.Next(ctx) {
		var row domain.OperatorSchema
		err := cursor.Decode(&row)
		helper.PanicIfError(err)
		result = append(result, row)
	}
	return result
}

func (repository *OperatorRepositoryImpl) ResetPasswordById(ctx context.Context, db *mongo.Database, operatorId string) {
	id, err := primitive.ObjectIDFromHex(operatorId)
	helper.PanicIfError(err)
	hash, err := helper.HashPassword("defaultpassword")
	helper.PanicIfError(err)
	res := db.Collection("student").FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"password": hash,
	}})
	if res.Err().Error() == "mongo: no documents in result" {
		panic(exception.NewNotFoundError("Data not found"))
	}
}

func (repository *OperatorRepositoryImpl) Destroy(ctx context.Context, db *mongo.Database, operatorId string) {
	id, err := primitive.ObjectIDFromHex(operatorId)
	helper.PanicIfError(err)
	res := db.Collection("student").FindOneAndDelete(ctx, bson.M{"_id": id})
	if res.Err().Error() == "mongo: no documents in result" {
		panic(exception.NewNotFoundError("Data not found"))
	} else {
		helper.PanicIfError(err)
	}
}
