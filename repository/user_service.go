package repository

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Login(ctx context.Context, db *mongo.Database, request domain.Login) string
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
