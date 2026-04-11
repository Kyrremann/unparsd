package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/jessevdk/go-flags"
	"github.com/kyrremann/unparsd/v4/fetch"
	"github.com/kyrremann/unparsd/v4/parsing"
	"github.com/kyrremann/unparsd/v4/statistics"
)

// generateCommand implements the 'generate' subcommand.
type generateCommand struct {
	Untappd   string `short:"u" long:"untappd" description:"Directory of per-year check-in JSON files" value-name:"DIR" default:"./checkins"`
	Output    string `short:"o" long:"output" description:"Output directory for generated statistics files" value-name:"DIR" default:"./"`
	AllStyles string `short:"s" long:"all-styles" description:"Path to all-styles.json; omit to scrape Untappd live" value-name:"FILE"`
	Username  string `short:"n" long:"username" description:"Namespace outputs under _data/<username>/ and _monthly/<username>/" value-name:"USERNAME"`
}

func (p *generateCommand) Execute(_ []string) error {
	if err := os.MkdirAll(p.Output, 0o750); err != nil {
		return err
	}

	if _, err := os.Stat(p.Untappd); err != nil {
		return fmt.Errorf("input path %q not found: %w", p.Untappd, err)
	}

	db, err := parsing.LoadJsonIntoDatabase(p.Untappd)
	if err != nil {
		return err
	}

	base, err := filepath.Abs(p.Output)
	if err != nil {
		return err
	}

	dataPath := filepath.Join(base, "_data")
	monthlyPath := filepath.Join(base, "_monthly")
	if p.Username != "" {
		dataPath = filepath.Join(dataPath, p.Username)
		monthlyPath = filepath.Join(monthlyPath, p.Username)
	}

	if err := os.MkdirAll(dataPath, 0o750); err != nil {
		return err
	}
	if err := os.MkdirAll(monthlyPath, 0o750); err != nil {
		return err
	}

	if err := statistics.GenerateAndSave(db, dataPath, p.AllStyles); err != nil {
		return err
	}

	return statistics.GenerateMonthlyAndSave(db, monthlyPath, p.Username)
}

// fetchCommand implements the 'fetch' subcommand.
type fetchCommand struct {
	Username string `short:"u" long:"username" description:"Untappd username to fetch check-ins for" required:"true"`
	Output   string `short:"o" long:"output" description:"Directory to write per-year JSON files into" default:"./checkins"`
}

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
	parser := flags.NewParser(nil, flags.Default)
	if _, err := parser.AddCommand(
		"generate",
		"Generate statistics from per-year check-in files",
		"Reads a directory of per-year JSON files (e.g. checkins/) and\n"+
			"writes statistics JSON files to the output directory.",
		&generateCommand{},
	); err != nil {
		log.Fatalf("registering generate command: %v", err)
	}

	if _, err := parser.AddCommand(
		"fetch",
		"Fetch check-ins from the Untappd API",
		"Downloads check-ins for a user and saves them as per-year JSON files.\n"+
			"Reads credentials from UNTAPPD_CLIENT_ID and UNTAPPD_CLIENT_SECRET.",
		&fetchCommand{},
	); err != nil {
		log.Fatalf("registering fetch command: %v", err)
	}

	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
}
