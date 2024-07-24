package login

import (
	"github.com/spf13/cobra"
)

const (
	loginExample = `  # Login interactively with a single sign-on provider
  omnistrate-ctl login

  # Login with email and password
  omnistrate-ctl login --email email --password password

  # Login with email and password from stdin. Save the password in a file and use cat to read it
  cat ~/omnistrate_pass.txt | omnistrate-ctl login --email email --password-stdin

  # Login with email and password from stdin. Save the password in an environment variable and use echo to read it
  echo $OMNISTRATE_PASSWORD | omnistrate-ctl login --email email --password-stdin`
)

var (
	email         string
	password      string
	passwordStdin bool
)

// LoginCmd represents the login command
var LoginCmd = &cobra.Command{
	Use:          `login`,
	Short:        "Log in to the Omnistrate platform.",
	Long:         `The login command is used to authenticate and log in to the Omnistrate platform.`,
	Example:      loginExample,
	RunE:         runLogin,
	SilenceUsage: true,
}

func init() {
	LoginCmd.Flags().StringVarP(&email, "email", "", "", "email")
	LoginCmd.Flags().StringVarP(&password, "password", "", "", "password")
	LoginCmd.Flags().BoolVarP(&passwordStdin, "password-stdin", "", false, "Reads the password from stdin")
}

func runLogin(cmd *cobra.Command, args []string) error {
	defer resetLogin()

	// Login with SSO if no email and password provided
	if len(email) == 0 && len(password) == 0 && !passwordStdin {
		return SSOLogin(cmd, args)
	}

	// Login with email and password
	return PasswordLogin(cmd, args)
}

func resetLogin() {
	email = ""
	password = ""
	passwordStdin = false
}
