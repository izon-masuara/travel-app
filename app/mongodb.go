package app

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	return client.Database(os.Getenv("DB_NAME"))
}
