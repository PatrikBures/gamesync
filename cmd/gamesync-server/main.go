package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("required mode, valid modes: serve")
	}

	mode := os.Args[1]
	switch mode {
	case "serve":
		if err := serve(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("invalid mode: '%s'", mode)
	}
}
