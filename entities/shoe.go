package entities

import "time"

type Shoe struct {
	Shoe_Id          uint
	Shoe_Name        string
	Shoe_Brand       Brand
	Shoe_Type        string
	Shoe_Description string
	Shoe_Sku         string
	Shoe_Price       int64
	Shoe_Stock       int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
