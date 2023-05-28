package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Search struct {
	IdDest primitive.ObjectID `bson:"id_dest" json:"id_dest"`
	Title  string             `bson:"title" json:"title"`
	Name   string             `bson:"name" json:"name"`
}
