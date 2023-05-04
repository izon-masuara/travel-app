package domain

import "time"

type DestinationUpdate struct {
	Title     string    `json:"title" bson:"title"`
	Date      time.Time `json:"date" bson:"date"`
	ImageFile string    `json:"image_file" bson:"image_file"`
	Location  Location  `json:"location" bson:"location"`
	Text      string    `json:"text" bson:"text"`
}
