package statistics

import (
	"github.com/kyrremann/unparsd/models"
	"gorm.io/gorm"
)

type Brewery struct {
	models.Brewery

	ListOfBeers string `json:"beers"`
	Checkins    int    `json:"checkins"`
}

func BreweryStats(db *gorm.DB) ([]Brewery, error) {
	var breweries []Brewery
	res := db.
		Table("breweries").
		Select("breweries.id as id," +
			"breweries.name," +
			"breweries.country," +
			"count(checkins.beer_id) as checkins").
		Joins("LEFT JOIN beers ON beers.brewery_id == breweries.id").
		Joins("LEFT JOIN checkins ON beers.id == checkins.beer_id").
		Group("beers.brewery_id").
		Find(&breweries)

	if res.Error != nil {
		return nil, res.Error
	}

	for i, brewery := range breweries {
		var beers string
		res = db.
			Table("beers").
			Distinct("name").
			Select("group_concat(name, '\n') as list_of_beers").
			Where("beers.brewery_id == ?", brewery.ID).
			Find(&beers)

		if res.Error != nil {
			return nil, res.Error
		}

		breweries[i].ListOfBeers = beers
	}

	return breweries, nil
}
