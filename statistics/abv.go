package statistics

import (
	"gorm.io/gorm"
)

// ABVBucket represents the check-in count within an ABV range.
type ABVBucket struct {
	Label string  `json:"label"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int     `json:"count"`
}

// ABVDistribution returns check-in counts grouped into ABV buckets:
// 0–3%, 3–5%, 5–7%, 7–10%, 10%+.
func ABVDistribution(db *gorm.DB) ([]ABVBucket, error) {
	buckets := []ABVBucket{
		{Label: "0-3%", Min: 0, Max: 3},
		{Label: "3-5%", Min: 3, Max: 5},
		{Label: "5-7%", Min: 5, Max: 7},
		{Label: "7-10%", Min: 7, Max: 10},
		{Label: "10%+", Min: 10, Max: 9999},
	}

	for i, b := range buckets {
		var count int64
		var res *gorm.DB
		if b.Max >= 9999 {
			res = db.Table("checkins").
				Joins("LEFT JOIN beers ON beers.id = checkins.beer_id").
				Where("beers.abv >= ?", b.Min).
				Count(&count)
		} else {
			res = db.Table("checkins").
				Joins("LEFT JOIN beers ON beers.id = checkins.beer_id").
				Where("beers.abv >= ? AND beers.abv < ?", b.Min, b.Max).
				Count(&count)
		}
		if res.Error != nil {
			return nil, res.Error
		}
		buckets[i].Count = int(count)
	}

	return buckets, nil
}
