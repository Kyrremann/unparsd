package statistics

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kyrremann/unparsd/models"
	"github.com/kyrremann/unparsd/parsing"
	"gorm.io/gorm"
)

type DistinctStyle struct {
	Type     string `json:"type"`
	Distinct int    `gorm:"-" json:"distinct"`
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

func getStylesFromUntappd() ([]string, error) {
	resp, err := http.Get("https://untappd.com/beer/top_rated")
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var styles []string
	doc.Find("select#filter_picker").Find("option").Each(func(i int, s *goquery.Selection) {
		style := strings.TrimSpace(s.Text())
		if style != "Show All Styles" {
			styles = append(styles, style)
		}
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
	var styles []DistinctStyle
	var distinctive []DistinctStyle
	res := db.Model(&models.Beer{}).Select("Type, count(Type) as total").Group("Type").Find(&distinctive)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, d := range distinctive {
		styles = append(styles, DistinctStyle{Type: d.Type, Distinct: d.Total})
	}

	var checkins []DistinctStyle
	res = db.Model(&models.Checkin{}).Select("checkins.beer").Joins("Beer").Select("Beer.Type, count(Beer.Type) as total").Group("Beer.Type").Find(&checkins)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, c := range checkins {
		for i, style := range styles {
			if style.Type == c.Type {
				style.Total = c.Total
				styles[i] = style
				break
			}
		}
	}

	return styles, nil
}
