package service

import (
	"context"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/entity/web"
	"kautsar/travel-app-api/helper"
	"kautsar/travel-app-api/repository"

	"go.mongodb.org/mongo-driver/mongo"
)

type OperatorService interface {
	Create(ctx context.Context, request *web.OperatorCreateRequest)
}

type OperatorServiceImpl struct {
	OperatorRepository repository.OperatorRepository
	Db                 *mongo.Database
}

func (service *OperatorServiceImpl) Create(ctx context.Context, request *web.OperatorCreateRequest) {
	operatorPayload := domain.Operator{
		Name:     request.Name,
		Username: request.Username,
		Password: request.Password,
	}
	err := service.OperatorRepository.Save(ctx, service.Db, &operatorPayload)
	helper.PanicIfError(err)
}
