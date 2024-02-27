/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/omnistrate/ctl/config"
	"github.com/spf13/cobra"
)

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:     "logout",
	Short:   "Logout from Omnistrate platform",
	Long:    `Logout from Omnistrate platform.`,
	Example: `omnistrate-cli logout`,
	RunE:    runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	err := config.RemoveAuthConfig()
	if err != nil {
		return err
	}
	fmt.Println("credentials removed")

	return nil
}
