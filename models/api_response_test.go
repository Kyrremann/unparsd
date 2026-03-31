package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeDate(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "Thu, 01 Mar 2016 19:06:42 +0000",
			want:  "2016-03-01 19:06:42",
		},
		{
			input: "Mon, 31 Dec 2018 23:59:59 +0000",
			want:  "2018-12-31 23:59:59",
		},
		// Non-RFC1123Z input passes through unchanged.
		{
			input: "2016-03-01 19:06:42",
			want:  "2016-03-01 19:06:42",
		},
	}
	for _, tc := range tests {
		got := normalizeDate(tc.input)
		assert.Equal(t, tc.want, got, "input: %s", tc.input)
	}
}

func TestJoinTaggedFriends(t *testing.T) {
	friends := []APITaggedFriend{
		{UID: 1, FirstName: "Alice", LastName: "Smith"},
		{UID: 2, FirstName: "Bob", LastName: "Jones"},
	}
	assert.Equal(t, "Alice Smith, Bob Jones", joinTaggedFriends(friends))
}

func TestJoinTaggedFriendsEmpty(t *testing.T) {
	assert.Equal(t, "", joinTaggedFriends(nil))
}

func TestJoinTaggedFriendsFallbackToUsername(t *testing.T) {
	friends := []APITaggedFriend{
		{UID: 42, UserName: "user42", FirstName: "", LastName: ""},
	}
	assert.Equal(t, "user42", joinTaggedFriends(friends))
}

func TestToJSONCheckin(t *testing.T) {
	api := &APICheckin{
		CheckinID:      99,
		CreatedAt:      "Thu, 01 Mar 2016 19:06:42 +0000",
		CheckinComment: "Great beer",
		RatingScore:    4.5,
		ServingType:    "Draft",
		FlavorProfiles: "Hoppy, Bitter",
		PurchaseVenue:  "Local store",
		Beer: APIBeer{
			BID:       1001,
			BeerName:  "Test IPA",
			BeerAbv:   6.5,
			BeerIbu:   60,
			BeerStyle: "IPA - American",
		},
		Brewery: APIBrewery{
			BreweryID:    5,
			BreweryName:  "Test Brewery",
			BreweryCity:  "Oslo",
			BreweryState: "",
			CountryName:  "Norway",
		},
		Venue: APIVenue{
			VenueID:      10,
			VenueName:    "The Pub",
			VenueCity:    "Oslo",
			VenueState:   "",
			VenueCountry: "Norway",
			Location:     APILocation{VenueLat: 59.91, VenueLng: 10.75},
		},
		Toasts:   APICount{Count: 3},
		Comments: APICount{Count: 1},
		Media: APIMedia{
			Count: 1,
			Items: []APIMediaItem{{Photo: APIPhoto{PhotoImgLg: "https://example.com/photo.jpg"}}},
		},
		TaggedFriends: APITaggedFriends{
			Count: 1,
			Items: []APITaggedFriend{{UID: 7, FirstName: "Eve", LastName: "Hansen"}},
		},
	}

	got := api.ToJSONCheckin()

	assert.Equal(t, 99, got.CheckinID)
	assert.Equal(t, "2016-03-01 19:06:42", got.CheckinAt)
	assert.Equal(t, "Great beer", got.Comment)
	assert.Equal(t, 1001, got.BID)
	assert.Equal(t, "Test IPA", got.BeerName)
	assert.Equal(t, float32(6.5), got.BeerAbv)
	assert.Equal(t, 60, got.BeerIbu)
	assert.Equal(t, "IPA - American", got.BeerType)
	assert.Equal(t, 5, got.BreweryID)
	assert.Equal(t, "Norway", got.BreweryCountry)
	assert.Equal(t, "The Pub", got.VenueName)
	assert.Equal(t, float32(59.91), got.VenueLat)
	assert.Equal(t, "Draft", got.ServingTypes)
	assert.Equal(t, "Hoppy, Bitter", got.FlavorProfiles)
	assert.Equal(t, "Local store", got.PurchaseVenue)
	assert.Equal(t, 3, got.TotalToasts)
	assert.Equal(t, 1, got.TotalComments)
	assert.Equal(t, "https://example.com/photo.jpg", got.PhotoUrl)
	assert.Equal(t, "Eve Hansen", got.TaggedFriends)
}

func TestToJSONCheckinNoMedia(t *testing.T) {
	api := &APICheckin{
		CheckinID: 1,
		CreatedAt: "Mon, 01 Jan 2018 00:00:00 +0000",
		Beer:      APIBeer{BID: 2},
	}
	got := api.ToJSONCheckin()
	assert.Equal(t, "", got.PhotoUrl)
	assert.Equal(t, "", got.TaggedFriends)
}
