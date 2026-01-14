package gamesyncstate

import (
	"encoding/json"
	"fmt"
	"gamesync/internal/state"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s DIR\n", os.Args[0])
	}

	dir := os.Args[1]

	state, err := state.GetState(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting state: %v\n", err)
	}

	stateJson, err := json.Marshal(state)

	fmt.Println(string(stateJson))
}
