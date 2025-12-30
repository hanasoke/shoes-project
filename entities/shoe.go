package entities

import "time"

type Shoe struct {
	Id          uint
	Name        string
	Brand       Brand
	Type        string
	Description string
	SKU         string
	Price       int64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ShoeCreate struct {
	Name        string
	IdBrand     uint
	Type        string
	Description string
	SKU         string
	Price       int64
	Stock       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ShoeUpdate struct {
	Name        string
	IdBrand     uint
	Type        string
	Description string
	SKU         string
	Price       int64
	Stock       int64
	UpdatedAt   time.Time
}
