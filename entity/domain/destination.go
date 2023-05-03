package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	Filename string `json:"filename" bson:"filename"`
	Mimetype string `json:"mimetype" bson:"mimetype"`
	Size     int    `json:"size" bson:"int"`
}

type Location struct {
	Long string `json:"long" bson:"long"`
	Lat  string `json:"lat" bson:"lat"`
}

type Facility struct {
	Name     string   `json:"name" bson:"name"`
	Location Location `json:"location" bson:"location"`
	Type     string   `json:"type" bson:"type"`
}

type Comment struct {
	Username string `json:"username" bson:"username"`
	Message  string `json:"message" bson:"message"`
}

type Destination struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	Title      string             `json:"title" bson:"title"`
	Date       time.Time          `json:"date" bson:"date"`
	ImageFile  Image              `json:"image_file" bson:"image_file"`
	Location   Location           `json:"location" bson:"location"`
	Text       string             `json:"text" bson:"text"`
	Rate       float64            `json:"rete" bson:"rate"`
	Facilities []Facility         `json:"facilities" bson:"facilities"`
	Comments   []Comment          `json:"comments" bson:"comments"`
}
