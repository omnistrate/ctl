package cmd

import (
	"github.com/omnistrate/ctl/config"
	"os"
)

func cleanup() {
	_ = os.RemoveAll(config.ConfigDir())
}
