package utils

import (
	"os"
	"strings"
)

func GetEnv(env, defaultValue string) string {
	environment := os.Getenv(env)
	if environment == "" {
		return defaultValue
	}
	return environment
}

func GetEnvAsBoolean(envName, defaultValue string) bool {
	return strings.ToLower(GetEnv(envName, defaultValue)) == "true"
}
