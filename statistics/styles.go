package statistics

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyrremann/unparsd/v4/models"
	"github.com/kyrremann/unparsd/v4/parsing"
	"gorm.io/gorm"
)

type DistinctStyle struct {
	Type     string `json:"type"`
	Distinct int    `gorm:"column:distinct_count" json:"distinct"`
	Total    int    `json:"total"`
}

func intersection(a, b []string) []string {
	var c []string

	for _, el := range a {
		if !slices.Contains(b, el) {
			c = append(c, el)
		}
	}

	return c
}

var excludedPrefixes = []string{"Wine", "Non-Alcoholic", "RTD", "THC Drink", "Spirits"}

func getStylesFromUntappd() (styles []string, err error) {
	resp, err := http.Get("https://untappd.com/beer/top_rated")
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	doc.Find("select#filter_picker").Find("option").Each(func(i int, s *goquery.Selection) {
		style := strings.TrimSpace(s.Text())
		if style == "Show All Styles" {
			return
		}

		lower := strings.ToLower(style)
		for _, prefix := range excludedPrefixes {
			if strings.HasPrefix(lower, strings.ToLower(prefix)) {
				return
			}
		}
		styles = append(styles, style)
	})

	return styles, nil
}

func getStylesDefinition(allStylesFile string) ([]string, error) {
	if len(allStylesFile) == 0 {
		return getStylesFromUntappd()
	}

	var styles []string
	err := parsing.ParseJsonFile(allStylesFile, &styles)
	return styles, err
}

func MissingStyles(db *gorm.DB, allStylesFile string) ([]string, error) {
	allStyles, err := getStylesDefinition(allStylesFile)
	if err != nil {
		return nil, err
	}

	var styles []string
	res := db.Model(&models.Beer{}).Distinct("Type").Pluck("Type", &styles)
	if res.Error != nil {
		return nil, res.Error
	}

	return intersection(allStyles, styles), nil
}

func DistinctStyles(db *gorm.DB) ([]DistinctStyle, error) {
	// One query: distinct beer count per style + total checkin count per style.
	var styles []DistinctStyle
	res := db.
		Table("beers").
		Select("beers.type as type," +
			"count(DISTINCT beers.id) as distinct_count," +
			"count(checkins.id) as total").
		Joins("LEFT JOIN checkins ON checkins.beer_id = beers.id").
		Group("beers.type").
		Find(&styles)
	return styles, res.Error
}
