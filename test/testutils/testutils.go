package testutils

import (
	"os"

	"github.com/pkg/errors"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
)

func Cleanup() {
	_ = os.RemoveAll(config.ConfigDir())
}

func Contains(arr []string, s string) bool {
	for _, a := range arr {
		if a == s {
			return true
		}
	}
	return false
}

func GetSmokeTestAccount() (string, string, error) {
	email := utils.GetEnv("SMOKE_TEST_EMAIL", "not-set")
	password := utils.GetEnv("SMOKE_TEST_PASSWORD", "")
	if email == "not-set" {
		return "", "", errors.New("SMOKE_TEST_EMAIL environment variable is not set. Set the environment variable to run the smoke test")
	}
	if password == "" {
		return "", "", errors.New("SMOKE_TEST_PASSWORD environment variable is not set. Set the environment variable to run the smoke test")
	}
	return email, password, nil
}
