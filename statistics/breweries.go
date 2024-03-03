package statistics

import (
	"sort"
	"strings"

	"github.com/kyrremann/unparsd/models"
	"github.com/pariz/gountries"
	"gorm.io/gorm"
)

type Brewery struct {
	models.Brewery

	ListOfBeers   string `json:"beers"`
	Checkins      int    `json:"checkins"`
	ISO3166Alpha2 string `json:"ISO3166Alpha2"`
}

func BreweryStats(db *gorm.DB) ([]Brewery, error) {
	var breweries []Brewery
	res := db.
		Table("breweries").
		Select("breweries.id as id," +
			"breweries.name," +
			"breweries.country," +
			"breweries.state," +
			"count(checkins.beer_id) as checkins").
		Joins("LEFT JOIN beers ON beers.brewery_id == breweries.id").
		Joins("LEFT JOIN checkins ON beers.id == checkins.beer_id").
		Group("beers.brewery_id").
		Find(&breweries)

	if res.Error != nil {
		return nil, res.Error
	}

	iso := ISO3166Alpha2{
		Query: gountries.New(),
	}

	for i, brewery := range breweries {
		var beers []string
		res = db.
			Table("beers").
			Distinct("name").
			Select("name").
			Where("beers.brewery_id == ?", brewery.ID).
			Order("beers.name ASC").
			Find(&beers)

		if res.Error != nil {
			return nil, res.Error
		}

		sort.Strings(beers)
		breweries[i].ListOfBeers = strings.Join(beers, "\n")

		_, ISO3166Alpha2, err := iso.getISO3166Alpha2(brewery.Country, brewery.State)
		if err != nil {
			return nil, err
		}
		breweries[i].ISO3166Alpha2 = ISO3166Alpha2
	}

	return breweries, nil
}
