package main

import (
	"log"

	"github.com/kyrremann/unparsd/parsing"
)

func main() {
	_, err := parsing.LoadJsonIntoDatabase("untappd.json")
	if err != nil {
		log.Fatal(err)
	}
}
