package statistics

import (
	"gorm.io/gorm"
)

// ServingTypeStat holds check-in count for a single serving type.
type ServingTypeStat struct {
	ServingType string `json:"serving_type"`
	Count       int    `json:"count"`
}

// ServingTypeStats returns check-in counts grouped by serving type,
// sorted by count descending. Check-ins without a serving type are excluded.
func ServingTypeStats(db *gorm.DB) ([]ServingTypeStat, error) {
	var stats []ServingTypeStat
	res := db.
		Table("checkins").
		Select("serving_types as serving_type, count(*) as count").
		Where("serving_types != ''").
		Group("serving_types").
		Order("count DESC").
		Find(&stats)
	return stats, res.Error
}
