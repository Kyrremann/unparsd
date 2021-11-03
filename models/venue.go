package models

type Venue struct {
	Name    string `gorm:"primaryKey"`
	City    string
	State   string
	Country string
	Lat     string
	Lng     string
}
