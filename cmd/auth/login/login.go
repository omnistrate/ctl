package login

import (
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/cqroot/prompt/input"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"regexp"
)

type loginMethod string

const (
	loginExample = `  # Login interactively with a single sign-on provider or using email and password
  omnistrate-ctl login

  # Login with email and password
  omnistrate-ctl login --email email --password password

  # Login with environment variables
  export OMNISTRATE_USER_NAME=YOUR_EMAIL
  export OMNISTRATE_PASSWORD=YOUR_PASSWORD
  ./omnistrate-ctl-darwin-arm64 login --email "$OMNISTRATE_USER_NAME" --password "$OMNISTRATE_PASSWORD"

  # Login with email and password from stdin. Save the password in a file and use cat to read it
  cat ~/omnistrate_pass.txt | omnistrate-ctl login --email email --password-stdin

  # Login with email and password from stdin. Save the password in an environment variable and use echo to read it
  echo $OMNISTRATE_PASSWORD | omnistrate-ctl login --email email --password-stdin`

	loginWithEmailAndPassword loginMethod = "Login with email and password"
	loginWithGoogle           loginMethod = "Login with Google"
	loginWithGitHub           loginMethod = "Login with GitHub"
)

var (
	email         string
	password      string
	passwordStdin bool
	gh            bool
	google        bool
)

// LoginCmd represents the login command
var LoginCmd = &cobra.Command{
	Use:          `login`,
	Short:        "Log in to the Omnistrate platform",
	Long:         `The login command is used to authenticate and log in to the Omnistrate platform.`,
	Example:      loginExample,
	RunE:         runLogin,
	SilenceUsage: true,
}

func init() {
	LoginCmd.Flags().StringVarP(&email, "email", "", "", "email")
	LoginCmd.Flags().StringVarP(&password, "password", "", "", "password")
	LoginCmd.Flags().BoolVarP(&passwordStdin, "password-stdin", "", false, "Reads the password from stdin")

	LoginCmd.Flags().BoolVarP(&gh, "gh", "", false, "Login with GitHub")
	LoginCmd.Flags().BoolVarP(&google, "google", "", false, "Login with Google")

	LoginCmd.MarkFlagsMutuallyExclusive("gh", "google", "email")

	LoginCmd.Args = cobra.NoArgs
}

func runLogin(cmd *cobra.Command, args []string) error {
	defer resetLogin()

	// Login with email and password if any of the flags are set
	if len(email) > 0 || len(password) > 0 || passwordStdin {
		return PasswordLogin(cmd, args, false)
	}

	if gh {
		return SSOLogin(identityProviderGitHub)
	}

	if google {
		return SSOLogin(identityProviderGoogle)
	}

	// Login interactively
	choice, err := prompt.New().Ask("How would you like to log in?").
		Choose([]string{
			string(loginWithEmailAndPassword),
			string(loginWithGoogle),
			string(loginWithGitHub),
		}, choose.WithTheme(choose.ThemeArrow))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	switch choice {
	case string(loginWithEmailAndPassword):
		email, err = prompt.New().Ask("Please enter your email:").
			Input("", input.WithValidateFunc(
				func(input string) error {
					emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
					if emailRegex.MatchString(input) {
						return nil
					} else {
						return errors.New("invalid email address")
					}
				}))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		password, err = prompt.New().Ask("Please enter your password:").
			Input("", input.WithEchoMode(input.EchoPassword))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		return PasswordLogin(cmd, args, true)
	case string(loginWithGoogle):
		return SSOLogin(identityProviderGoogle)
	case string(loginWithGitHub):
		return SSOLogin(identityProviderGitHub)

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
