package models

import (
	"fmt"
	"strings"
	"time"
)

// APIResponse is the top-level envelope returned by the Untappd API.
type APIResponse struct {
	Meta     APIMeta            `json:"meta"`
	Response APICheckinResponse `json:"response"`
}

type APIMeta struct {
	Code int `json:"code"`
}

type APICheckinResponse struct {
	Checkins APICheckinList `json:"checkins"`
}

type APICheckinList struct {
	Count      int           `json:"count"`
	Items      []APICheckin  `json:"items"`
	Pagination APIPagination `json:"pagination"`
}

type APIPagination struct {
	MaxID int `json:"max_id"`
}

// APICheckin is a single check-in item from the Untappd API response.
// The nested structure differs from the flat JSONCheckin export format.
type APICheckin struct {
	CheckinID      int              `json:"checkin_id"`
	CreatedAt      string           `json:"created_at"`
	CheckinComment string           `json:"checkin_comment"`
	RatingScore    any              `json:"rating_score"`
	ServingType    string           `json:"serving_type"`
	FlavorProfiles string           `json:"flavor_profiles"`
	PurchaseVenue  string           `json:"purchase_venue"`
	Beer           APIBeer          `json:"beer"`
	Brewery        APIBrewery       `json:"brewery"`
	Venue          APIVenue         `json:"venue"`
	Toasts         APICount         `json:"toasts"`
	Comments       APICount         `json:"comments"`
	Media          APIMedia         `json:"media"`
	TaggedFriends  APITaggedFriends `json:"tagged_friends"`
}

type APIBeer struct {
	BID       int     `json:"bid"`
	BeerName  string  `json:"beer_name"`
	BeerAbv   float32 `json:"beer_abv"`
	BeerIbu   int     `json:"beer_ibu"`
	BeerStyle string  `json:"beer_style"`
}

type APIBrewery struct {
	BreweryID    int    `json:"brewery_id"`
	BreweryName  string `json:"brewery_name"`
	BreweryCity  string `json:"brewery_city"`
	BreweryState string `json:"brewery_state"`
	CountryName  string `json:"country_name"`
}

type APIVenue struct {
	VenueID      int         `json:"venue_id"`
	VenueName    string      `json:"venue_name"`
	VenueCity    string      `json:"venue_city"`
	VenueState   string      `json:"venue_state"`
	VenueCountry string      `json:"venue_country"`
	Location     APILocation `json:"location"`
}

type APILocation struct {
	VenueLat float32 `json:"venue_lat"`
	VenueLng float32 `json:"venue_lng"`
}

type APICount struct {
	Count int `json:"count"`
}

type APIMedia struct {
	Count int            `json:"count"`
	Items []APIMediaItem `json:"items"`
}

type APIMediaItem struct {
	Photo APIPhoto `json:"photo"`
}

type APIPhoto struct {
	PhotoImgLg string `json:"photo_img_lg"`
}

type APITaggedFriends struct {
	Count int               `json:"count"`
	Items []APITaggedFriend `json:"items"`
}

type APITaggedFriend struct {
	UID       int    `json:"uid"`
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// normalizeDate converts the Untappd API date format (RFC1123Z, e.g.
// "Thu, 01 Mar 2016 19:06:42 +0000") to the flat export format used
// throughout this codebase ("2006-01-02 15:04:05").
func normalizeDate(apiDate string) string {
	t, err := time.Parse(time.RFC1123Z, apiDate)
	if err != nil {
		return apiDate
	}
	return t.UTC().Format("2006-01-02 15:04:05")
}

// joinTaggedFriends converts the API array of tagged friends to a
// comma-separated "First Last" string matching the export format.
func joinTaggedFriends(friends []APITaggedFriend) string {
	names := make([]string, 0, len(friends))
	for _, f := range friends {
		name := strings.TrimSpace(fmt.Sprintf("%s %s", f.FirstName, f.LastName))
		if name == "" {
			name = f.UserName
		}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

// ToJSONCheckin converts a nested APICheckin into the canonical flat
// JSONCheckin format used for storage and statistics generation.
// Fields not available via the API (GlobalRatingScore,
// GlobalWeightedRatingScore) are left at their zero values.
func (a *APICheckin) ToJSONCheckin() JSONCheckin {
	photoURL := ""
	if len(a.Media.Items) > 0 {
		photoURL = a.Media.Items[0].Photo.PhotoImgLg
	}

	return JSONCheckin{
		BID:            a.Beer.BID,
		BeerAbv:        a.Beer.BeerAbv,
		BeerIbu:        a.Beer.BeerIbu,
		BeerName:       a.Beer.BeerName,
		BeerType:       a.Beer.BeerStyle,
		BreweryCity:    a.Brewery.BreweryCity,
		BreweryCountry: a.Brewery.CountryName,
		BreweryID:      a.Brewery.BreweryID,
		BreweryName:    a.Brewery.BreweryName,
		BreweryState:   a.Brewery.BreweryState,
		CheckinAt:      normalizeDate(a.CreatedAt),
		CheckinID:      a.CheckinID,
		Comment:        a.CheckinComment,
		FlavorProfiles: a.FlavorProfiles,
		PhotoUrl:       photoURL,
		PurchaseVenue:  a.PurchaseVenue,
		RatingScore:    a.RatingScore,
		ServingTypes:   a.ServingType,
		TaggedFriends:  joinTaggedFriends(a.TaggedFriends.Items),
		TotalComments:  a.Comments.Count,
		TotalToasts:    a.Toasts.Count,
		VenueCity:      a.Venue.VenueCity,
		VenueCountry:   a.Venue.VenueCountry,
		VenueLat:       a.Venue.Location.VenueLat,
		VenueLng:       a.Venue.Location.VenueLng,
		VenueName:      a.Venue.VenueName,
		VenueState:     a.Venue.VenueState,
	}
}
