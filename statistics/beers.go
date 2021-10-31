package statistics

import (
	"gorm.io/gorm"
)

type Beer struct {
	Name      string
	AvgRating float64
	CheckinAt string
	Type      string
	Checkins  int
	Ibu       int
	Abv       float32
	Brewery   string
}

func BeerStats(db *gorm.DB) ([]Beer, error) {
	var beers []Beer
	res := db.
		Table("checkins").
		Select("beers.name as name, AVG(checkins.rating_score) as avg_rating," +
			"count(checkins.id) as checkins," +
			"checkins.checkin_at," +
			"beers.ibu, beers.abv," +
			"breweries.name as brewery").
		Joins("INNER JOIN beers ON beers.id = checkins.beer_id").
		Joins("INNER JOIN breweries ON breweries.id = beers.brewery_id").
		Group("checkins.beer_id").
		Find(&beers)

	if res.Error != nil {
		return nil, res.Error
	}

	return beers, nil
}
