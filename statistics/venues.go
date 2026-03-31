package statistics

import (
	"gorm.io/gorm"
)

// VenueStat holds check-in count and metadata for a single venue.
type VenueStat struct {
	Name     string  `json:"name"`
	City     string  `json:"city"`
	State    string  `json:"state"`
	Country  string  `json:"country"`
	Lat      float32 `json:"lat"`
	Lng      float32 `json:"lng"`
	Checkins int     `json:"checkins"`
}

// TopVenues returns all named venues sorted by check-in count descending.
func TopVenues(db *gorm.DB) ([]VenueStat, error) {
	var stats []VenueStat
	res := db.
		Table("checkins").
		Select("checkins.venue_name as name," +
			"venues.city," +
			"venues.state," +
			"venues.country," +
			"venues.lat," +
			"venues.lng," +
			"count(checkins.id) as checkins").
		Joins("LEFT JOIN venues ON venues.name = checkins.venue_name").
		Where("checkins.venue_name != ''").
		Group("checkins.venue_name").
		Order("checkins DESC").
		Find(&stats)
	return stats, res.Error
}
