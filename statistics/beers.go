package statistics

import (
	"github.com/kyrremann/unparsd/models"
	"gorm.io/gorm"
)

type Beer struct {
	models.Beer

	AvgRating   float64 `json:"avg_rating"`
	CheckinAt   string  `json:"checkin_at"`
	Checkins    int     `json:"checkins"`
	BreweryName string  `json:"brewery"`
}

func BeerStats(db *gorm.DB) ([]Beer, error) {
	var beers []Beer
	res := db.
		Table("checkins").
		Select("beers.id as id," +
			"beers.name as name," +
			"ROUND(AVG(checkins.rating_score), 2) as avg_rating," +
			"count(checkins.id) as checkins," +
			"beers.type as type," +
			"strftime('%Y-%m-%d', checkins.checkin_at) as checkin_at," +
			"beers.ibu, beers.abv," +
			"breweries.name as brewery_name," +
			"breweries.id as brewery_id").
		Joins("LEFT JOIN beers ON beers.id = checkins.beer_id").
		Joins("LEFT JOIN breweries ON breweries.id = beers.brewery_id").
		Group("checkins.beer_id").
		Find(&beers)

	if res.Error != nil {
		return nil, res.Error
	}

	return beers, nil
}
