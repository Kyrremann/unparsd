package main

import (
	_ "embed"
	"path/filepath"

	"github.com/jessevdk/go-flags"
	"github.com/kyrremann/unparsd/parsing"
	"github.com/kyrremann/unparsd/statistics"
)

//go:embed fixture/all_styles.json
var allStylesJson []byte

var opts struct {
	Untappd string `short:"u" long:"untappd" description:"" value-name:"untappd.json" default:"untappd.json"`
	Output  string `short:"o" long:"output" description:"" value-name:"_data" default:"_data"`
}

func main() {
	statistics.AllStylesJson = allStylesJson

	_, err := flags.Parse(&opts)
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

	err = statistics.GenerateAndSave(db, base)
	if err != nil {
		panic(err)
	}
}
