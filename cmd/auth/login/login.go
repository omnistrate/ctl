package login

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	ctlutils "github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	goa "goa.design/goa/v3/pkg"
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

	// SSO login if no email and password provided
	if len(email) == 0 && len(password) == 0 && !passwordStdin {
		return SsoLogin(cmd, args)
	}

	// Otherwise, login with email and password
	if len(password) > 0 {
		ctlutils.PrintWarning("WARNING! Using --password is insecure, consider using the --password-stdin flag. Check the help for examples.")
		if passwordStdin {
			err := fmt.Errorf("--password and --password-stdin are mutually exclusive")
			ctlutils.PrintError(err)
			return err
		}

		if len(email) == 0 {
			err := errors.New("must provide --email with --password")
			ctlutils.PrintError(err)
			return err
		}
	}

	if passwordStdin {
		if len(email) == 0 {
			err := errors.New("must provide --email with --password-stdin")
			ctlutils.PrintError(err)
			return err
		}

		passwordFromStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			ctlutils.PrintError(err)
			return err
		}
		password = strings.TrimSpace(string(passwordFromStdin))
	}

	password = strings.TrimSpace(password)
	if len(password) == 0 {
		err := errors.New("must provide a non-empty password via --password or --password-stdin")
		ctlutils.PrintError(err)
		return err
	}

	ctlutils.PrintInfo("Calling the Omnistrate server to validate the credentials...")

	token, err := validateLogin(email, password)
	if err != nil {
		ctlutils.PrintError(err)
		return err
	}

	authConfig := config.AuthConfig{
		Token: token,
	}
	if err = config.CreateOrUpdateAuthConfig(authConfig); err != nil {
		ctlutils.PrintError(err)
		return err
	}

	ctlutils.PrintSuccess(fmt.Sprintf("Credential saved for %s", email))

	return nil
}

func validateLogin(email string, pass string) (token string, err error) {
	signin, err := httpclientwrapper.NewSignin(ctlutils.GetHostScheme(), ctlutils.GetHost())
	if err != nil {
		return "", err
	}

	request := signinapi.SigninRequest{
		Email:    email,
		Password: utils.ToPtr(pass),
	}

	res, err := signin.Signin(context.Background(), &request)
	if err != nil {
		var serviceErr *goa.ServiceError
		ok := errors.As(err, &serviceErr)
		if !ok {
			return
		}

		return "", fmt.Errorf("%s\nDetail: %s", serviceErr.Name, serviceErr.Message)
	}

	token = res.JWTToken
	return
}

func resetLogin() {
	email = ""
	password = ""
	passwordStdin = false
}
