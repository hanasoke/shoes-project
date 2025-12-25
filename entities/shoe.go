package entities

import "time"

type Shoe struct {
	Id          uint
	Name        string
	Brand       Brand
	Type        string
	Description string
	Sku         string
	Price       int64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
