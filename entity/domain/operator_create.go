package domain

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type OperatorCreate struct {
	Name     string         `json:"name" bson:"name"`
	Username string         `json:"username" bson:"username"`
	Role     string         `json:"role"`
	Password string         `json:"password" bson:"password"`
	Position [][]Coordinate `json:"position" bson:"position"`
}
