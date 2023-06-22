package service

import (
	"context"
	"database/sql"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/repository"
)

type TiketService interface {
	FindOne(ctx context.Context, request string) domain.Tiket
}

type TiketServiceImpl struct {
	TiketRepository repository.TiketRepository
	Db              *sql.DB
}

func NewTiketService(tiketRepository repository.TiketRepository, db *sql.DB) TiketService {
	return &TiketServiceImpl{
		TiketRepository: tiketRepository,
		Db:              db,
	}
}

func (service TiketServiceImpl) FindOne(ctx context.Context, request string) domain.Tiket {
	tiket := service.TiketRepository.FindOne(ctx, service.Db, request)
	return tiket
}
