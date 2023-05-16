package service

import (
	"context"
	"fmt"
	"io/ioutil"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/repository"
	"os"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	Login(ctx context.Context, requset web.LoginRequest) string
	FindDestinationByRegion(ctx context.Context, requset string) []domain.Destination
	FindOneDestinationByRegion(ctx context.Context, requestRegion string, requestDestination string) domain.Destination
	FindAllRegions(ctx context.Context) []string
	FindImage(imgName string) []byte
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

func (service *UserServiceImpl) FindDestinationByRegion(ctx context.Context, requset string) []domain.Destination {
	result := service.UserRepository.FindDestinationByRegion(ctx, service.Db, requset)
	return result
}

func (service *UserServiceImpl) FindOneDestinationByRegion(ctx context.Context, requestRegion string, requestDestination string) domain.Destination {
	result := service.UserRepository.FindOneDestinationByRegion(ctx, service.Db, requestRegion, requestDestination)
	return result
}

func (service *UserServiceImpl) FindAllRegions(ctx context.Context) []string {
	regions := service.UserRepository.FindAllRegions(ctx, service.Db)
	return regions
}

func (service *UserServiceImpl) FindImage(imgName string) []byte {
	dir, err := os.Getwd()
	helper.PanicIfError(err)
	path := fmt.Sprintf("%s/public/%s", dir, imgName)
	fileByte, err := ioutil.ReadFile(path)
	helper.PanicIfError(err)
	return fileByte
}
