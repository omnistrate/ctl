package testutils

import (
	"errors"
	"os"

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
		return "", "", errors.New("TEST_EMAIL is not set")
	}
	if password == "" {
		return "", "", errors.New("TEST_PASSWORD is not set")
	}
	return email, password, nil
}
