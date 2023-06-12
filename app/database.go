package app

import (
	"context"
	"kautsar/travel-app-api/helper"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewDb() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.TODO())
	helper.PanicIfError(err)

	return client.Database(os.Getenv("DB_NAME"))
}
