package models

// gorm.Model definition
type Brewery struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"-"`
	State   string `json:"-"`
	Country string `json:"country"`
	Beers   []Beer `json:"-"`
}

func (b *Brewery) TableName() string {
	return "breweries"
}
