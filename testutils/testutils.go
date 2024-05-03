package testutils

import (
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	"os"
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

// GetTestAccount returns the test account email and password for the test environment
func GetTestAccount() (string, string) {
	switch utils.GetEnv("ROOT_DOMAIN", "omnistrate.cloud") {
	case "omnistrate.dev":
		return "xzhang+customer-hosted@omnistrate.com", "Test@1234"
	default:
		return "", ""
	}
}
