package web

import "time"

type DestinationUpdateRequest struct {
	Title     string    `validate:"required" json:"title" form:"title"`
	Date      time.Time `validate:"required" json:"date" form:"date"`
	Long      string    `validate:"required" json:"long" form:"long"`
	Lat       string    `validate:"required" json:"lat" form:"lat"`
	ImageFile string    `validate:"required" json:"image_file" form:"image_file"`
	Text      string    `validate:"required" json:"text" form:"text"`
}
