package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "omnistrate-ctl",
	Short: "Manage your Omnistrate SaaS from the command line.",
	Long: wordwrap.WrapString(`
Omnistrate ctl is a command line tool designed to streamline the creation, deployment, and management of your Omnistrate SaaS. Use it to build services from docker-compose files, manage service plans, and interact with the Omnistrate platform efficiently.

For additional support, please refer to the CTL reference documentation at https://docs.omnistrate.com/getting-started/ctl-reference/.`, 80),
	Run: runRoot,
}

func runRoot(cmd *cobra.Command, args []string) {
	// Check if the version flag is set
	versionFlag, err := cmd.Flags().GetBool("version")
	if err == nil && versionFlag {
		fmt.Println(Version)
		return
	}

	printLogo()
	err = cmd.Help()
	if err != nil {
		return
	}
}

// printLogo prints an ASCII logo, which was generated with figlet
func printLogo() {
	fmt.Println()
	colors := []color.Attribute{
		color.FgRed, color.FgYellow, color.FgGreen, color.FgCyan, color.FgBlue, color.FgMagenta,
	}
	for i, r := range figletStr {
		fmt.Printf("%s", color.New(colors[i%len(colors)]).SprintFunc()(string(r)))
	}
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
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version number of omnistrate-ctl")
}
