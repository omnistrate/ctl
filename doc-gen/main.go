package main

import (
	"log"
	"os"

	"github.com/omnistrate/ctl/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	// Set the KUBECONFIG environment variable to an empty string to avoid loading any kubeconfig settings of the user.
	os.Setenv("KUBECONFIG", "")
	err := doc.GenMarkdownTree(cmd.RootCmd, "./mkdocs/docs/")
	if err != nil {
		log.Fatal(err)
	}
}
