package domain

type OperatorCreate struct {
	Name     string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Role     string `json:"role"`
	Password string `json:"password" bson:"password"`
}
