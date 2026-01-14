package main

import (
	"encoding/json"
	"fmt"
	"gamesync/internal/state"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s DIR\n", os.Args[0])
		os.Exit(10)
	}

	dir := os.Args[1]

	state, err := state.Get(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting state: %v\n", err)
		os.Exit(11)
	}

	stateJson, err := json.Marshal(state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to json: %v\n", err)
		os.Exit(12)
	}

	fmt.Println(string(stateJson))
}
