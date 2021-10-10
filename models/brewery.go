package models

// gorm.Model definition
type Brewery struct {
	ID      int
	Name    string
	City    string
	State   string
	Country string
	Beers   []Beer
}

func (b *Brewery) TableName() string {
	return "breweries"
}
