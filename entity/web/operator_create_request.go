package web

type OperatorCreateRequest struct {
	Name     string `validate:"required,min=4,max=50" json:"name" bson:"name"`
	Username string `validate:"required,min=4,max=50" json:"username" bson:"username"`
	Password string `validate:"required,min=8,max=50" json:"password" bson:"password"`
}
