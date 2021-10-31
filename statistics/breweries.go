package statistics

import "gorm.io/gorm"

type Brewery struct {
	Name     string
	Country  string
	State    string
	City     string
	Beers    string
	Checkins int
}

func BreweryStats(db *gorm.DB) ([]Brewery, error) {
	var breweries []Brewery
	res := db.
		Table("breweries").
		Select("breweries.name, country, city, state," +
			"group_concat(beers.name) as beers," +
			"count(checkins.beer_id) as checkins").
		Joins("INNER JOIN beers ON beers.brewery_id == breweries.id").
		Joins("INNER JOIN checkins ON beers.id == checkins.beer_id").
		Group("beers.brewery_id").
		Find(&breweries)

	if res.Error != nil {
		return nil, res.Error
	}

	return breweries, nil
}
