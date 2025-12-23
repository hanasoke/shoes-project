package entities

import "time"

type Shoe struct {
	Shoe_Id          uint
	Shoe_name        string
	brand            Brand
	shoe_type        string
	shoe_size        string
	show_description string
	shoe_sku         string
	shoe_price       int64
	shoe_stock       int64
	shoe_photo       string
	shoe_status      string
	created_at       time.Time
	updated_at       time.Time
}
