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

type UserRepository interface {
	Login(ctx context.Context, db *mongo.Database, request domain.Login) string
	FindDestinationByRegion(ctx context.Context, db *mongo.Database, request string) []domain.Destination
	FindOneDestinationByRegion(ctx context.Context, db *mongo.Database, requestRegion string, requestDestination string) domain.Destination
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Login(ctx context.Context, db *mongo.Database, request domain.Login) string {
	var found domain.Account
	err := db.Collection("account").FindOne(ctx, bson.M{
		"username": request.Username,
	}).Decode(&found)
	if err != nil {
		panic(exception.NewAuthError("Username or password are wrong"))
	}
	if !helper.DecribePassword(request.Passowrd, found.Password) {
		panic(exception.NewAuthError("Username or password are wrong"))
	}
	token := helper.GenerateToken(&helper.JwtPayload{
		Name: found.Name,
		Role: found.Role,
	})
	return token
}

func (repository *UserRepositoryImpl) FindDestinationByRegion(ctx context.Context, db *mongo.Database, request string) []domain.Destination {
	csr, err := db.Collection(request).Find(ctx, bson.D{})
	helper.PanicIfError(err)
	defer csr.Close(ctx)
	var result []domain.Destination
	for csr.Next(ctx) {
		var row domain.Destination
		err := csr.Decode(&row)
		helper.PanicIfError(err)
		result = append(result, row)
	}
	return result
}

func (repository *UserRepositoryImpl) FindOneDestinationByRegion(ctx context.Context, db *mongo.Database, requestRegion string, requestDestination string) domain.Destination {
	id, err := primitive.ObjectIDFromHex(requestDestination)
	helper.PanicIfError(err)
	var found domain.Destination
	err = db.Collection(requestRegion).FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&found)
	if err != nil {
		panic(exception.NewNotFoundError("Data not found"))
	}
	return found
}
