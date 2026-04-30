package main

import (
	"gamesync/internal/api"
	"log"
)

func main() {
	if err := api.Serve(); err != nil {
		log.Fatal(err)
	}
}
