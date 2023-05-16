package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type OperatorSchema struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Username string             `json:"username" bson:"username"`
}
