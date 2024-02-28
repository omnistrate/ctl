package testutils

import (
	"github.com/omnistrate/ctl/config"
	"os"
)

func Cleanup() {
	_ = os.RemoveAll(config.ConfigDir())
}
