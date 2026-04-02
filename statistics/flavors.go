package statistics

import (
	"sort"
	"strings"

	"gorm.io/gorm"
)

// FlavorStat holds the frequency of a single flavor profile tag.
type FlavorStat struct {
	Flavor string `json:"flavor"`
	Count  int    `json:"count"`
}

// FlavorProfileStats splits the comma-separated flavor_profiles field in Go
// and returns a frequency table sorted by count descending.
func FlavorProfileStats(db *gorm.DB) ([]FlavorStat, error) {
	// Pull all non-empty flavor_profiles strings.
	var profiles []string
	res := db.
		Table("checkins").
		Where("flavor_profiles != ''").
		Pluck("flavor_profiles", &profiles)
	if res.Error != nil {
		return nil, res.Error
	}

	freq := make(map[string]int)
	for _, p := range profiles {
		for tag := range strings.SplitSeq(p, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				freq[tag]++
			}
		}
	}

	stats := make([]FlavorStat, 0, len(freq))
	for flavor, count := range freq {
		stats = append(stats, FlavorStat{Flavor: flavor, Count: count})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Count != stats[j].Count {
			return stats[i].Count > stats[j].Count
		}
		return stats[i].Flavor < stats[j].Flavor
	})

	return stats, nil
}
