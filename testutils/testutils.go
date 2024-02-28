package testutils

import (
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
