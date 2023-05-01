package repository

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OperatorRepository interface {
	Save(ctx context.Context, db *mongo.Database, operator domain.OperatorCreate) error
	FindAll(ctx context.Context, db *mongo.Database) []domain.OperatorSchema
}

type OperatorRepositoryImpl struct {
}

func NewOperatorRepository() OperatorRepository {
	return &OperatorRepositoryImpl{}
}

// func connect() (*mongo.Database, error) {
// 	clientOptions := options.Client()
// 	clientOptions.ApplyURI("mongodb://localhost:27017")
// 	client, err := mongo.NewClient(clientOptions)
// 	helper.PanicIfError(err)
// 	err = client.Connect(context.Background())
// 	helper.PanicIfError(err)

// 	return client.Database("belajar_golang"), nil
// }

func (repository *OperatorRepositoryImpl) Save(ctx context.Context, db *mongo.Database, operator domain.OperatorCreate) error {
	_, err := db.Collection("student").InsertOne(ctx, operator)
	helper.PanicIfError(err)
	return nil
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
