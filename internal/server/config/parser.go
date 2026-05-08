package serverConfig

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)


func AddBoolVar(variable *bool, name string, defaultValue bool, description string) {
	flag.BoolVar(variable, name, defaultValue, description)

	if os.Getenv(getEnvName(name)) == "true" {
		*variable = true
	}
}

func AddStringVar(variable *string, name string, defaultValue string, description string) {
	flag.StringVar(variable, name, defaultValue, description)

	if e := os.Getenv(getEnvName(name)); e != "" {
		*variable = e
	}
}

func AddIntvar(variable *int, name string, defaultValue int, description string) {
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
