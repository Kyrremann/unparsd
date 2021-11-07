package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/kyrremann/unparsd/parsing"
	"github.com/kyrremann/unparsd/statistics"
)

var opts struct {
	Untappd   string `short:"u" long:"untappd" description:"" value-name:"untappd.json" default:"untappd.json"`
	Output    string `short:"o" long:"output" description:"" value-name:"_data" default:"./"`
	AllStyles string `short:"s" long:"all-styles" description:"" value-name:"all-styles.json"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(opts.Output, 0755)
	if err != nil {
		panic(err)
	}
	_, err = os.Stat(opts.Untappd)
	if err != nil {
		panic(err)
	}

	db, err := parsing.LoadJsonIntoDatabase(opts.Untappd)
	if err != nil {
		panic(err)
	}

	base, err := filepath.Abs(opts.Output)
	if err != nil {
		panic(err)
	}

	err = statistics.GenerateAndSave(db, base, opts.AllStyles)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}
}
