package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Login struct {
	Username string `json:"username" bson:"username"`
	Passowrd string `json:"password" bson:"password"`
}

type Account struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
}
