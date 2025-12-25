package entities

import "time"

type Brand struct {
	Id        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
