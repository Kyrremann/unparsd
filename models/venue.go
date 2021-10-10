package models

import "gorm.io/gorm"

type Venue struct {
	gorm.Model

	Name    string `gorm:"primaryKey"`
	City    string
	State   string
	Country string
	Lat     string
	Lng     string
}
