package login

import (
	"github.com/burl/inquire"
	"github.com/burl/inquire/widget"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	loginExample = `  # Login interactively with a single sign-on provider or using email and password
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

	// Login with email and password if any of the flags are set
	if len(email) > 0 || len(password) > 0 || passwordStdin {
		return PasswordLogin(cmd, args)
	}

	// Login interactively
	var identityProvider string
	inquire.Query().
		Menu(&identityProvider, "How would you like to log in? [Use arrows to move, enter to select]", func(w *widget.Menu) {
			w.Item("Password", "Login with email and password")
			w.Item("Google", "Login with Google")
			w.Item("GitHub", "Login with GitHub")
		}).
		Exec()

	switch identityProvider {
	case "Password":
		inquire.Query().Input(&email, "Please enter your email", func(w *widget.Input) {
			w.Valid(func(input string) string {
				if !strings.Contains(input, "@") {
					return "Please enter a valid email address"
				}
				return ""
			})
		}).
			Input(&password, "Please enter your password", func(w *widget.Input) {
				w.MaskInput()
			}).
			Exec()

		return PasswordLogin(cmd, args)
	case "Google":
		return SSOLoginWithGoogle(cmd, args)
	case "GitHub":
		return SSOLoginWithGitHub(cmd, args)

	default:
		err := errors.New("Invalid selection")
		utils.PrintError(err)
		return err
	}
}

func resetLogin() {
	email = ""
	password = ""
	passwordStdin = false
}
