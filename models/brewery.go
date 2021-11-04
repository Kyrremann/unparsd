package models

// gorm.Model definition
type Brewery struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	City    string `json:"country"`
	State   string `json:"-"`
	Country string `json:"-"`
	Beers   []Beer `json:"-"`
}

func (b *Brewery) TableName() string {
	return "breweries"
}
