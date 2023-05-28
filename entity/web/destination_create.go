package web

import (
	"kautsar/travel-app-api/entity/domain"
	"time"
)

type DestinationCreateRequest struct {
	Title      string            `validate:"required" form:"title"`
	Date       time.Time         `validate:"required" form:"date"`
	Long       string            `validate:"required" form:"long"`
	Lat        string            `validate:"required" form:"lat"`
	ImageFile  string            `validate:"required" form:"image_file"`
	Text       string            `validate:"required" form:"text"`
	Facilities []domain.Facility `validate:"required"`
}
