/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omnistrate-ctl",
	Short: "Omnistrate ctl in Go.",
	Long: `Omnistrate ctl is a command line tool for creating,
	deploying, and managing your Omnistrate SaaS. `,
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	printLogo()
	cmd.Help()
}

// printLogo prints an ASCII logo, which was generated with figlet
func printLogo() {
	fmt.Printf(figletStr)
}

const figletStr = `                  _     __           __     
 ___  __ _  ___  (_)__ / /________ _/ /____ 
/ _ \/  ' \/ _ \/ (_-</ __/ __/ _ ` + "`" + `/ __/ -_)
\___/_/_/_/_//_/_/___/\__/_/  \_,_/\__/\__/ 

`

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ctl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
