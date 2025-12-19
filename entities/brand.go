package entities

import "time"

type Brand struct {
	Brand_Id   uint
	Brand_Name string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
