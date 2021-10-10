package models

// gorm.Model definition
type Beer struct {
	ID                        int
	Name                      string
	Type                      string
	Abv                       float32
	Ibu                       int
	GlobalWeightedRatingScore float32
	GlobalRatingScore         float32
	BreweryID                 int
	Brewery                   Brewery
}
