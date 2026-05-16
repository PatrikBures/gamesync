package serverConfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)


const sliceSeperator = '|'

// Adds a bool flag and env var.
// The env var value must be "true" for the variable to be true.
//
// The env var is turned uppercase and replaces "-" to "_" and
// adds "GAMESYNC_" prefix to it.
func AddBoolVar(variable *bool, name string, defaultValue bool, description string) {
	flag.BoolVar(variable, name, defaultValue, description)

	if os.Getenv(getEnvName(name)) == "true" {
		*variable = true
	}
}

// Adds a string flag and env var. 
//
// The env var is turned uppercase and replaces "-" to "_" and
// adds "GAMESYNC_" prefix to it.
func AddStringVar(variable *string, name string, defaultValue string, description string) {
	flag.StringVar(variable, name, defaultValue, description)

	if e := os.Getenv(getEnvName(name)); e != "" {
		*variable = e
	}
}

// Adds a int flag and env var.
// If the env var is set it needs to be a valid int, otherwise it panics.
//
// The env var is turned uppercase and replaces "-" to "_" and
// adds "GAMESYNC_" prefix to it.
func AddIntVar(variable *int, name string, defaultValue int, description string) {
	flag.IntVar(variable, name, defaultValue, description)

	envName := getEnvName(name)
	sInt := os.Getenv(envName)
	if sInt == "" {
		return
	}
	i, err := strconv.Atoi(sInt)
	if err != nil {
		panic(fmt.Errorf("env var '%s' has invalid int: %v", envName, err))
	}
	*variable = i
}

func getEnvName(name string) string {
	return "GAMESYNC_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}

func StringToSlice(s string) []string {
	return strings.Split(s, string(sliceSeperator))
}
