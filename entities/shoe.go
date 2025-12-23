package entities

import "time"

type Shoe struct {
	Id          uint
	Name        string
	Brand       Brand
	Type        string
	Size        string
	Description string
	Sku         string
	Price       int64
	Stock       int64
	Photo       string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
