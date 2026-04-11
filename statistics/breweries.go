package statistics

import (
	"sort"
	"strings"

	"github.com/kyrremann/unparsd/v4/models"
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

	// Pre-load all beer names grouped by brewery in one query.
	type beerRow struct {
		BreweryID int
		Name      string
	}
	var beerRows []beerRow
	res = db.Table("beers").
		Select("brewery_id, name").
		Distinct("brewery_id, name").
		Order("name ASC").
		Find(&beerRows)
	if res.Error != nil {
		return nil, res.Error
	}

	beersByBrewery := make(map[int][]string, len(beerRows))
	for _, br := range beerRows {
		beersByBrewery[br.BreweryID] = append(beersByBrewery[br.BreweryID], br.Name)
	}

	for i, brewery := range breweries {
		names := beersByBrewery[brewery.ID]
		sort.Strings(names)
		breweries[i].ListOfBeers = strings.Join(names, "\n")

		_, ISO3166Alpha2, err := getISO3166Alpha2(brewery.Country, brewery.State)
		if err != nil {
			return nil, err
		}
		breweries[i].ISO3166Alpha2 = ISO3166Alpha2
	}

	return breweries, nil
}
