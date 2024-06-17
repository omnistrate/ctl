package cmd

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

var (
	email         string
	password      string
	passwordStdin bool
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   `login [--email EMAIL] [--password PASSWORD]`,
	Short: "Log in to the Omnistrate platform",
	Long:  `The login command is used to authenticate and log in to the Omnistrate platform.`,
	Example: `  cat ~/omnistrate_pass.txt | omnistrate-ctl login --email email --password-stdin
	  echo $OMNISTRATE_PASSWORD | omnistrate-ctl login --email email --password-stdin`,
	RunE:         runLogin,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&email, "email", "", "", "email")
	loginCmd.Flags().StringVarP(&password, "password", "", "", "password")
	loginCmd.Flags().BoolVarP(&passwordStdin, "password-stdin", "", false, "Reads the password from stdin")
}

func runLogin(cmd *cobra.Command, args []string) error {
	defer resetLogin()

	if len(email) == 0 {
		err := errors.New("must provide --email")
		ctlutils.PrintError(err)
		return err
	}

	if len(password) > 0 {
		ctlutils.PrintWarning("WARNING! Using --password is insecure, consider using: cat ~/omnistrate_pass.txt | omnistrate-ctl login -e email --password-stdin echo $PASSWORD")
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
		Email: email,
		Token: token,
		Auth:  config.JWTAuthType,
	}
	if err = config.UpdateAuthConfig(authConfig); err != nil {
		ctlutils.PrintError(err)
		return err
	}

	authConfig, err = config.LookupAuthConfig()
	if err != nil {
		ctlutils.PrintError(err)
		return err
	}

	ctlutils.PrintSuccess(fmt.Sprintf("Credential saved for %s", authConfig.Email))

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
