package service

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/repository"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	Login(ctx context.Context, requset web.LoginRequest) string
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository
	Db             *mongo.Database
	Validate       *validator.Validate
}

func NewUserService(userRespository repository.UserRepository, db *mongo.Database, validate *validator.Validate) UserService {
	return &UserServiceImpl{
		UserRepository: userRespository,
		Db:             db,
		Validate:       validate,
	}
}

func (service *UserServiceImpl) Login(ctx context.Context, request web.LoginRequest) string {
	err := service.Validate.Struct(request)
	helper.PanicIfError(err)

	loginPayload := domain.Login{
		Username: request.Username,
		Passowrd: request.Password,
	}

	token := service.UserRepository.Login(ctx, service.Db, loginPayload)
	return token
}
