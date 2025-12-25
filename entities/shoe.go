package entities

import "time"

type Shoe struct {
	ID          uint
	Name        string
	Brand_Name  string
	Type        string
	Description string
	SKU         string
	Price       int64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
