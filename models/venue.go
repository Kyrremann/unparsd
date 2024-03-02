package models

type Venue struct {
	Name    string `gorm:"primaryKey"`
	City    string
	State   string
	Country string
	Lat     float32
	Lng     float32
}
