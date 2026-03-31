package statistics

import (
	"math"

	"gorm.io/gorm"
)

// RatingDelta holds a beer with the delta between the user's average rating
// and the global weighted rating. Sorted by absolute delta descending.
type RatingDelta struct {
	BeerID       int     `json:"beer_id"`
	BeerName     string  `json:"beer_name"`
	BreweryName  string  `json:"brewery_name"`
	AvgPersonal  float64 `json:"avg_personal"`
	GlobalRating float64 `json:"global_rating"`
	Delta        float64 `json:"delta"`
	Checkins     int     `json:"checkins"`
}

// RatingDeltas returns all beers that have at least one rated check-in,
// along with the delta between personal and global rating, sorted by
// absolute delta descending.
func RatingDeltas(db *gorm.DB) ([]RatingDelta, error) {
	type row struct {
		BeerID               int     `gorm:"column:beer_id"`
		BeerName             string  `gorm:"column:beer_name"`
		BreweryName          string  `gorm:"column:brewery_name"`
		AvgPersonal          float64 `gorm:"column:avg_personal"`
		GlobalWeightedRating float64 `gorm:"column:global_weighted_rating"`
		Checkins             int     `gorm:"column:checkins"`
	}

	var rows []row
	res := db.
		Table("checkins").
		Select("checkins.beer_id," +
			"beers.name as beer_name," +
			"breweries.name as brewery_name," +
			"ROUND(AVG(checkins.rating_score), 2) as avg_personal," +
			"beers.global_weighted_rating_score as global_weighted_rating," +
			"count(checkins.id) as checkins").
		Joins("LEFT JOIN beers ON beers.id = checkins.beer_id").
		Joins("LEFT JOIN breweries ON breweries.id = beers.brewery_id").
		Where("checkins.rating_score > 0").
		Group("checkins.beer_id").
		Find(&rows)
	if res.Error != nil {
		return nil, res.Error
	}

	deltas := make([]RatingDelta, 0, len(rows))
	for _, r := range rows {
		delta := math.Round((r.AvgPersonal-r.GlobalWeightedRating)*100) / 100
		deltas = append(deltas, RatingDelta{
			BeerID:       r.BeerID,
			BeerName:     r.BeerName,
			BreweryName:  r.BreweryName,
			AvgPersonal:  r.AvgPersonal,
			GlobalRating: r.GlobalWeightedRating,
			Delta:        delta,
			Checkins:     r.Checkins,
		})
	}

	// Sort by absolute delta descending.
	for i := 0; i < len(deltas)-1; i++ {
		for j := i + 1; j < len(deltas); j++ {
			if math.Abs(deltas[j].Delta) > math.Abs(deltas[i].Delta) {
				deltas[i], deltas[j] = deltas[j], deltas[i]
			}
		}
	}

	return deltas, nil
}
