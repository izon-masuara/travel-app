package repository

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/mongo"
)

type OperatorRepository interface {
	Save(ctx context.Context, db *mongo.Database, operator *domain.Operator) error
	// Update(ctx context.Context, db *mongo.Database, operator *domain.Operator) error
	// Delete(ctx context.Context, db *mongo.Database, operator *domain.Operator)
	// FindById(ctx context.Context, db *mongo.Database, operatorId int) *domain.Operator
	// FindAll(ctx context.Context, db *mongo.Database) []*domain.Operator
}

type OperatorRepositoryImpl struct {
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

func (repository *OperatorRepositoryImpl) Save(ctx context.Context, db *mongo.Database, operator *domain.Operator) error {
	_, err := db.Collection("student").InsertOne(ctx, operator)
	helper.PanicIfError(err)
	return nil
}

// func (repository *OperatorRepositoryImpl) Update(ctx context.Context, db *mongo.Database, operator domain.Operator) error {
// 	db.Collection("student")
// 	return nil
// }
