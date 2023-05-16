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

type OperatorService interface {
	Create(ctx context.Context, request web.OperatorCreateRequest) string
	FindAll(ctx context.Context) []domain.OperatorSchema
	ResetPassword(ctx context.Context, request string) string
	Destroy(ctx context.Context, request string) string
}

type OperatorServiceImpl struct {
	OperatorRepository repository.OperatorRepository
	Db                 *mongo.Database
	validate           *validator.Validate
}

func NewOperatorService(operatorRepository repository.OperatorRepository, Db *mongo.Database, validate *validator.Validate) OperatorService {
	return &OperatorServiceImpl{
		OperatorRepository: operatorRepository,
		Db:                 Db,
		validate:           validate,
	}
}

func (service *OperatorServiceImpl) Create(ctx context.Context, request web.OperatorCreateRequest) string {
	err := service.validate.Struct(request)
	helper.PanicIfError(err)

	operatorPayload := domain.OperatorCreate{
		Name:     request.Name,
		Username: request.Username,
		Password: request.Password,
	}

	service.OperatorRepository.Save(ctx, service.Db, operatorPayload)
	return "Success create new account"
}

func (service *OperatorServiceImpl) FindAll(ctx context.Context) []domain.OperatorSchema {
	result := service.OperatorRepository.FindAll(ctx, service.Db)
	return result
}

func (service *OperatorServiceImpl) ResetPassword(ctx context.Context, request string) string {
	service.OperatorRepository.ResetPasswordById(ctx, service.Db, request)
	return "Success reset password"
}

func (service *OperatorServiceImpl) Destroy(ctx context.Context, request string) string {
	service.OperatorRepository.Destroy(ctx, service.Db, request)
	return "Success delete account"
}
