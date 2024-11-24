package util

import (
	"log"
	"os"
	"strconv"
)

func GetPositiveIntegerEnvironmentVariable(name string, defaultValue int, isValidFunction func() bool) int {
	rawVariable, found := os.LookupEnv(name)
	if !found {
		return defaultValue
	}

	variable, err := strconv.Atoi(rawVariable)
	if err != nil || variable < 0 {
		log.Fatalf("Positive integer required for %s", name)
	}

	if !isValidFunction() {
		log.Fatalf("Environment Variable %s is invalid", name)
	}

	return variable
}

func GetStringEnvironmentVariable(name string, defaultValue string, isValidFunction func() bool) string {
	variable, found := os.LookupEnv(name)

	if !found {
		variable = defaultValue
	}

	if !isValidFunction() {
		log.Fatalf("Environment Variable %s is invalid", name)
	}

	return variable
}
