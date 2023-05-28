package main

import (
	"context"
	"kautsar/travel-app-api/helper"
	"log"
	"os"
	"time"

	"github.com/subosito/gotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupTestDb() *mongo.Database {
	clientOptions := options.Client()
	clientOptions.ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.NewClient(clientOptions)
	helper.PanicIfError(err)
	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	helper.PanicIfError(err)

	return client.Database("travel_app_dev")
}

func startDb(db *mongo.Database) {
	err := db.Collection("account").Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	pass, _ := helper.HashPassword("admin123")
	_, err = db.Collection("account").InsertOne(context.Background(), bson.M{
		"name":     "admin",
		"username": "admin",
		"password": pass,
		"role":     "admin",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	err := gotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	db := setupTestDb()
	startDb(db)
}
