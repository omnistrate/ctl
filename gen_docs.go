package main

import (
	"log"

	"github.com/spf13/cobra/doc"
)

func main() {
	var myCmd = &cobra.Command{
		Use:   "mycmd",
		Short: "A brief description of your command",
		Long:  `A more detailed description of what your command does.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Your command logic here
		},
	}

	err := doc.GenMarkdownTree(myCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
