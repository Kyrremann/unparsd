package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/kyrremann/unparsd/fetch"
	"github.com/kyrremann/unparsd/parsing"
	"github.com/kyrremann/unparsd/statistics"
)

// Top-level options drive the default generate behaviour (backward-compatible).
var opts struct {
	Untappd   string `short:"u" long:"untappd" description:"Path to untappd.json file or directory of per-year JSON files" value-name:"untappd.json" default:"untappd.json"`
	Output    string `short:"o" long:"output" description:"Output directory for generated statistics files" value-name:"_data" default:"./"`
	AllStyles string `short:"s" long:"all-styles" description:"Path to all-styles.json; omit to scrape Untappd live" value-name:"all-styles.json"`
}

// fetchCommand implements the 'fetch' subcommand.
type fetchCommand struct {
	Username string `long:"username" description:"Untappd username to fetch check-ins for" required:"true"`
	Output   string `short:"o" long:"output" description:"Directory to write per-year JSON files into" default:"./checkins"`
}

// Execute is called automatically by go-flags when the 'fetch' subcommand
// is selected.  Credentials are read from environment variables so they
// never appear in shell history.
func (f *fetchCommand) Execute(_ []string) error {
	clientID := os.Getenv("UNTAPPD_CLIENT_ID")
	clientSecret := os.Getenv("UNTAPPD_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf(
			"UNTAPPD_CLIENT_ID and UNTAPPD_CLIENT_SECRET must be set\n" +
				"Register your app at https://untappd.com/api/register to obtain credentials",
		)
	}

	// Create a context that is cancelled on Ctrl+C or SIGTERM so that any
	// check-ins fetched so far are saved to disk before the process exits.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return fetch.FetchAndSave(ctx, f.Username, clientID, clientSecret, f.Output)
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.AddCommand(
		"fetch",
		"Fetch check-ins from the Untappd API",
		"Downloads check-ins for a user and saves them as per-year JSON files.\n"+
			"Reads credentials from UNTAPPD_CLIENT_ID and UNTAPPD_CLIENT_SECRET.",
		&fetchCommand{},
	)

	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	// A subcommand's Execute() already ran; nothing left to do.
	if parser.Active != nil {
		return
	}

	// No subcommand given – run the default generate flow.
	if err := generate(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func generate() error {
	if err := os.MkdirAll(opts.Output, 0755); err != nil {
		return err
	}

	if _, err := os.Stat(opts.Untappd); err != nil {
		return fmt.Errorf("input path %q not found: %w", opts.Untappd, err)
	}

	db, err := parsing.LoadJsonIntoDatabase(opts.Untappd)
	if err != nil {
		return err
	}

	base, err := filepath.Abs(opts.Output)
	if err != nil {
		return err
	}

	if err := statistics.GenerateAndSave(db, base, opts.AllStyles); err != nil {
		return err
	}

	return statistics.GenerateMonthlyAndSave(db, base)
}
