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

type DestinationService interface {
	Create(ctx context.Context, request web.DestinationCreateRequest) string
	FindAll(ctx context.Context) []domain.Destination
}

type DestinationServiceImpl struct {
	DestinationRepository repository.DestinationRepository
	Db                    *mongo.Database
	validate              *validator.Validate
}

func NewDestinationService(destinationRepository repository.DestinationRepository, Db *mongo.Database, validate *validator.Validate) DestinationService {
	return &DestinationServiceImpl{
		DestinationRepository: destinationRepository,
		Db:                    Db,
		validate:              validate,
	}
}

func (service *DestinationServiceImpl) Create(ctx context.Context, request web.DestinationCreateRequest) string {

	err := service.validate.Struct(request)
	if err != nil {
		helper.RemoveFile(request.ImageFile)
		helper.PanicIfError(err)
	}

	destinationPayload := domain.DestinationCreate{
		Title: request.Title,
		Date:  request.Date,
		Location: domain.Location{
			Long: request.Long,
			Lat:  request.Lat,
		},
		ImageFile:  request.ImageFile,
		Text:       request.Text,
		Rate:       0,
		Facilities: make([]domain.Facility, 0),
		Comments:   make([]domain.Comment, 0),
	}

	service.DestinationRepository.Save(ctx, service.Db, destinationPayload)
	return "Success create destination"
}

func (service *DestinationServiceImpl) FindAll(ctx context.Context) []domain.Destination {
	result := service.DestinationRepository.FindAll(ctx, service.Db)
	return result
}
