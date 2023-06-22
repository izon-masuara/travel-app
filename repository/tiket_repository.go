package repository

import (
	"context"
	"database/sql"
	"kautsar/travel-app-api/entity/domain"
	"kautsar/travel-app-api/exception"
	"kautsar/travel-app-api/helper"
)

type TiketRepository interface {
	FindOne(ctx context.Context, db *sql.DB, code string) domain.Tiket
}

type TiketRespositoryImpl struct {
}

func NewTiketRepository() TiketRepository {
	return &TiketRespositoryImpl{}
}

func (repository *TiketRespositoryImpl) FindOne(ctx context.Context, db *sql.DB, code string) domain.Tiket {
	SQL := "SELECT code, route, price, updated_at FROM tikets WHERE code = ?;"
	rows, err := db.QueryContext(ctx, SQL, code)
	helper.PanicIfError(err)
	tiket := domain.Tiket{}
	if rows.Next() {
		rows.Scan(
			&tiket.Code,
			&tiket.Route,
			&tiket.Price,
			&tiket.UpdatedAt,
		)
		return tiket
	} else {
		panic(exception.NewNotFoundError("Data not found"))
	}
}
