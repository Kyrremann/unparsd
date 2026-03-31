package fetch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/kyrremann/unparsd/models"
	"github.com/kyrremann/unparsd/parsing"
)

const (
	apiBaseURL = "https://api.untappd.com/v4"
	pageLimit  = 25
)

// FetchAndSave fetches all new check-ins for the given username from the
// Untappd API and saves them as per-year JSON files inside outputDir
// (e.g. outputDir/2023.json, outputDir/2024.json).
//
// On subsequent runs it loads existing files first and only downloads
// check-ins newer than the highest checkin_id already stored.
//
// If ctx is cancelled the function saves whatever has been fetched so far
// and returns nil so that no data is lost.
func FetchAndSave(ctx context.Context, username, clientID, clientSecret, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Load every check-in we already have on disk.
	existing, err := loadExistingCheckins(outputDir)
	if err != nil {
		return fmt.Errorf("loading existing check-ins: %w", err)
	}

	existingIDs := make(map[int]struct{}, len(existing))
	for _, c := range existing {
		existingIDs[c.CheckinID] = struct{}{}
	}
	log.Printf("Found %d existing check-ins", len(existing))

	// Fetch only check-ins that are not yet stored.
	newCheckins, fetchErr := fetchAllNew(ctx, username, clientID, clientSecret, existingIDs)

	// On cancellation, partial results are still worth saving.
	cancelled := errors.Is(fetchErr, context.Canceled) || errors.Is(fetchErr, context.DeadlineExceeded)
	if fetchErr != nil && !cancelled {
		return fmt.Errorf("fetching check-ins from API: %w", fetchErr)
	}

	if len(newCheckins) == 0 {
		if cancelled {
			log.Println("Fetch cancelled before any new check-ins were retrieved")
			return nil
		}
		log.Println("No new check-ins found")
		return nil
	}

	if cancelled {
		log.Printf("Fetch cancelled; saving %d partial check-ins", len(newCheckins))
	} else {
		log.Printf("Fetched %d new check-ins", len(newCheckins))
	}

	// Merge and split by year.
	return saveByYear(append(existing, newCheckins...), outputDir)
}

// loadExistingCheckins reads all *.json files from dir and returns the
// merged slice of check-ins.  Missing directory is treated as empty.
func loadExistingCheckins(dir string) ([]models.JSONCheckin, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var all []models.JSONCheckin
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		var checkins []models.JSONCheckin
		if err := parsing.ParseJsonFile(path, &checkins); err != nil {
			return nil, fmt.Errorf("reading %s: %w", path, err)
		}
		all = append(all, checkins...)
	}
	return all, nil
}

// fetchAllNew pages through the API newest-first, stopping as soon as it
// encounters a checkin_id that already exists in existingIDs.
// On context cancellation it returns whatever it has collected so far along
// with the context error so the caller can decide to save partial results.
func fetchAllNew(ctx context.Context, username, clientID, clientSecret string, existingIDs map[int]struct{}) ([]models.JSONCheckin, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	var newCheckins []models.JSONCheckin
	maxID := 0

	for {
		// Check for cancellation before each page fetch.
		if err := ctx.Err(); err != nil {
			return newCheckins, err
		}

		url := buildURL(username, clientID, clientSecret, maxID)
		items, nextMaxID, err := fetchPage(ctx, client, url)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return newCheckins, err
			}
			return nil, err
		}

		done := false
		for _, item := range items {
			checkin := item.ToJSONCheckin()
			if _, seen := existingIDs[checkin.CheckinID]; seen {
				done = true
				break
			}
			newCheckins = append(newCheckins, checkin)
		}

		if done || len(items) == 0 || nextMaxID == 0 {
			break
		}
		maxID = nextMaxID

		// Brief pause between pages to be a good API citizen.
		// Immediately interruptible on cancellation.
		select {
		case <-ctx.Done():
			return newCheckins, ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
	}

	return newCheckins, nil
}

func buildURL(username, clientID, clientSecret string, maxID int) string {
	url := fmt.Sprintf(
		"%s/user/checkins/%s?client_id=%s&client_secret=%s&limit=%d",
		apiBaseURL, username, clientID, clientSecret, pageLimit,
	)
	if maxID > 0 {
		url += fmt.Sprintf("&max_id=%d", maxID)
	}
	return url
}

func fetchPage(ctx context.Context, client *http.Client, url string) ([]models.APICheckin, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "unparsd (github.com/kyrremann/unparsd)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Honour rate-limit headers – pause when running low.
	// The sleep is interruptible via the context.
	if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining != "" {
		if n, err := strconv.Atoi(remaining); err == nil && n < 5 {
			log.Printf("Rate limit nearly exhausted (%d remaining), waiting 60s...", n)
			select {
			case <-ctx.Done():
				return nil, 0, ctx.Err()
			case <-time.After(60 * time.Second):
			}
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, 0, fmt.Errorf("parsing API response: %w", err)
	}

	items := apiResp.Response.Checkins.Items
	nextMaxID := apiResp.Response.Checkins.Pagination.MaxID
	// The Untappd API often returns max_id=0 in the pagination object.
	// Derive the next page cursor from the oldest item in the current page
	// instead (items are newest-first, so last item is oldest).
	if nextMaxID == 0 && len(items) > 0 {
		nextMaxID = items[len(items)-1].CheckinID
	}
	return items, nextMaxID, nil
}

// yearFromCheckin extracts the 4-digit year string from a check-in's
// CheckinAt field ("2006-01-02 15:04:05").
func yearFromCheckin(c models.JSONCheckin) string {
	t, err := time.Parse("2006-01-02 15:04:05", c.CheckinAt)
	if err != nil {
		if len(c.CheckinAt) >= 4 {
			return c.CheckinAt[:4]
		}
		return "unknown"
	}
	return strconv.Itoa(t.Year())
}

// saveByYear groups check-ins by year, sorts each group by checkin_id
// descending (newest first), and writes one JSON file per year to outputDir.
func saveByYear(checkins []models.JSONCheckin, outputDir string) error {
	byYear := make(map[string][]models.JSONCheckin)
	for _, c := range checkins {
		year := yearFromCheckin(c)
		byYear[year] = append(byYear[year], c)
	}

	for year, ycheckins := range byYear {
		sort.Slice(ycheckins, func(i, j int) bool {
			return ycheckins[i].CheckinID > ycheckins[j].CheckinID
		})

		path := filepath.Join(outputDir, year+".json")
		if err := parsing.SaveDataToJsonFile(ycheckins, path); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		log.Printf("Wrote %d check-ins to %s", len(ycheckins), path)
	}
	return nil
}
