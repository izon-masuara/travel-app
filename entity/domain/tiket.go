package domain

import "time"

type Tiket struct {
	Code      string    `json:"code"`
	Route     string    `json:"route"`
	Price     string    `json:"price"`
	UpdatedAt time.Time `json:"updated_at"`
}
