package app

import (
	"context"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDb() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	helper.PanicIfError(err)

	err = client.Connect(context.TODO())
	helper.PanicIfError(err)

	return client.Database("travel_app_dev")
}
