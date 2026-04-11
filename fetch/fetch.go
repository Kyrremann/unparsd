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

	"github.com/kyrremann/unparsd/v4/models"
	"github.com/kyrremann/unparsd/v4/parsing"
)

const (
	apiBaseURL = "https://api.untappd.com/v4"
	pageLimit  = 25
)

// FetchAndSave fetches all new check-ins for the given username from the
// Untappd API and saves them as per-year JSON files inside outputDir
// (e.g. outputDir/2023.json, outputDir/2024.json).
//
// On subsequent runs it reads the latest stored checkin_id, then pages
// backwards from today (newest-first) until it encounters that ID.
//
// If ctx is cancelled the function saves whatever has been fetched so far
// and returns nil so that no data is lost.
func FetchAndSave(ctx context.Context, username, clientID, clientSecret, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0o750); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	latestID, err := loadLatestCheckinID(outputDir)
	if err != nil {
		return fmt.Errorf("loading latest check-in ID: %w", err)
	}

	if latestID > 0 {
		log.Printf("Latest stored check-in ID: %d", latestID)
	} else {
		log.Println("No existing check-ins found; fetching all")
	}

	// Fetch only check-ins newer than latestID.
	newCheckins, fetchErr := fetchAllNew(ctx, username, clientID, clientSecret, latestID)

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

	return saveByYear(newCheckins, outputDir)
}

// loadLatestCheckinID reads the first item from each *.json file in dir
// (files are stored newest-first) and returns the highest checkin_id found.
// Returns 0 if no files exist.
func loadLatestCheckinID(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}

		return 0, err
	}

	latestID := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		id, err := readFirstCheckinID(path)
		if err != nil {
			return 0, fmt.Errorf("reading %s: %w", path, err)
		}

		if id > latestID {
			latestID = id
		}
	}

	return latestID, nil
}

// readFirstCheckinID opens a JSON file containing a []JSONCheckin array and
// decodes only the first element, returning its CheckinID.
func readFirstCheckinID(path string) (int, error) {
	// #nosec G304 -- path is constructed from outputDir provided by the caller (CLI flag)
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	// Consume the opening '['.
	if _, err := dec.Token(); err != nil {
		return 0, nil // empty or invalid file
	}

	var c models.JSONCheckin
	if err := dec.Decode(&c); err != nil {
		return 0, nil // empty array
	}
	return c.CheckinID, nil
}

// fetchAllNew pages backwards from the newest check-in, collecting every
// item with CheckinID > latestID.  It stops as soon as it encounters an ID
// that is already stored, or when a partial page signals no more history.
//
// On context cancellation it returns whatever it has collected so far.
func fetchAllNew(ctx context.Context, username, clientID, clientSecret string, latestID int) ([]models.JSONCheckin, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	var newCheckins []models.JSONCheckin
	cursor := 0 // 0 means start from the newest check-in (no max_id)

	for {
		// Check for cancellation before each page fetch.
		if err := ctx.Err(); err != nil {
			return newCheckins, err
		}

		url := buildURL(username, clientID, clientSecret, cursor)
		items, err := fetchPage(ctx, client, url)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return newCheckins, err
			}
			return nil, err
		}

		if len(items) == 0 {
			break
		}

		// Walk the page (newest → oldest); stop at the first ID we already have.
		done := false
		for _, item := range items {
			if item.CheckinID <= latestID {
				done = true
				break
			}

			newCheckins = append(newCheckins, item.ToJSONCheckin())
		}

		if done {
			break
		}

		// A partial page means there is no more history.
		if len(items) < pageLimit {
			break
		}

		// Advance the cursor to the oldest item on this page so the next
		// request returns the next batch going further back in time.
		cursor = items[len(items)-1].CheckinID

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

// buildURL constructs the API endpoint URL.
// cursor (max_id) paginates backwards through history; 0 means start from
// the newest check-in.
func buildURL(username, clientID, clientSecret string, cursor int) string {
	url := fmt.Sprintf(
		"%s/user/checkins/%s?client_id=%s&client_secret=%s&limit=%d",
		apiBaseURL, username, clientID, clientSecret, pageLimit,
	)

	if cursor > 0 {
		url += fmt.Sprintf("&max_id=%d", cursor)
	}

	return url
}

func fetchPage(ctx context.Context, client *http.Client, url string) ([]models.APICheckin, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "unparsd (github.com/kyrremann/unparsd)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Honour rate-limit headers – pause when running low.
	// The sleep is interruptible via the context.
	if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining != "" {
		if n, err := strconv.Atoi(remaining); err == nil && n < 5 {
			log.Printf("Rate limit nearly exhausted (%d remaining), waiting 60s...", n)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(60 * time.Second):
			}
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing API response: %w", err)
	}

	return apiResp.Response.Checkins.Items, nil
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

// saveByYear groups new check-ins by year, then for each affected year loads
// the existing file (if any), merges, sorts newest-first, and rewrites it.
// Only the year files that actually receive new check-ins are touched.
func saveByYear(newCheckins []models.JSONCheckin, outputDir string) error {
	byYear := make(map[string][]models.JSONCheckin)
	for _, c := range newCheckins {
		year := yearFromCheckin(c)
		byYear[year] = append(byYear[year], c)
	}

	for year, ycheckins := range byYear {
		path := filepath.Join(outputDir, year+".json")

		// Load the existing year file if it exists, then merge.
		var existing []models.JSONCheckin
		if _, err := os.Stat(path); err == nil {
			if err := parsing.ParseJsonFile(path, &existing); err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}
		}

		merged := append(existing, ycheckins...)
		sort.Slice(merged, func(i, j int) bool {
			return merged[i].CheckinID > merged[j].CheckinID
		})

		if err := parsing.SaveDataToJsonFile(merged, path); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		log.Printf("Wrote %d check-ins to %s", len(merged), path)
	}
	return nil
}
