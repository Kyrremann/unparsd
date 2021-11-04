package models

// gorm.Model definition
type Beer struct {
	ID                        int     `json:"id"`
	Name                      string  `json:"name"`
	Type                      string  `json:"type"`
	Abv                       float32 `json:"abv"`
	Ibu                       int     `json:"ibu"`
	GlobalWeightedRatingScore float32 `json:"-"`
	GlobalRatingScore         float32 `json:"-"`
	BreweryID                 int     `json:"brewery_id"`
	Brewery                   Brewery `json:"-"`
}
