package statistics

import (
	"gorm.io/gorm"
)

type Beer struct {
	Name      string  `json:"name"`
	AvgRating float64 `json:"avg_rating"`
	CheckinAt string  `json:"checkin_at"`
	Type      string  `json:"type"`
	Checkins  int     `json:"checkins"`
	Ibu       int     `json:"ibu"`
	Abv       float32 `json:"abv"`
	Brewery   string  `json:"brewery"`
}

func BeerStats(db *gorm.DB) ([]Beer, error) {
	var beers []Beer
	res := db.
		Table("checkins").
		Select("beers.name as name, AVG(checkins.rating_score) as avg_rating," +
			"count(checkins.id) as checkins," +
			"beers.type as type," +
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
