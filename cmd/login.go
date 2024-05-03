package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	utils2 "github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	goa "goa.design/goa/v3/pkg"
	"io"
	"os"
	"strings"
)

var (
	email         string
	password      string
	passwordStdin bool
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   `login [--email EMAIL] [--password PASSWORD]`,
	Short: "Log in to Omnistrate platform",
	Long:  "Log in to Omnistrate platform",
	Example: `  cat ~/omnistrate_pass.txt | ./omnistrate-cli login -e email --password-stdin
	  echo $OMNISTRATE_PASSWORD | ./omnistrate-cli login -e email --password-stdin`,
	RunE: runLogin,
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
		return fmt.Errorf("must provide --email or -e")
	}

	if len(password) > 0 {
		fmt.Println("WARNING! Using --password is insecure, consider using: cat ~/omnistrate_pass.txt | ./omnistrate-cli login -e email --password-stdin echo $PASSWORD")
		if passwordStdin {
			return fmt.Errorf("--password and --password-stdin are mutually exclusive")
		}

		if len(email) == 0 {
			return fmt.Errorf("must provide --email with --password")
		}
	}

	if passwordStdin {
		if len(email) == 0 {
			return fmt.Errorf("must provide --email with --password-stdin")
		}

		passwordFromStdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		password = strings.TrimSpace(string(passwordFromStdin))
	}

	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return fmt.Errorf("must provide a non-empty password via --password or --password-stdin")
	}

	fmt.Println("Calling the Omnistrate server to validate the credentials...")

	token, err := validateLogin(email, password)
	if err != nil {
		return err
	}

	authConfig := config.AuthConfig{
		Email: email,
		Token: token,
		Auth:  config.JWTAuthType,
	}
	if err = config.UpdateAuthConfig(authConfig); err != nil {
		return err
	}

	authConfig, err = config.LookupAuthConfig()
	if err != nil {
		return err
	}

	fmt.Println("credential saved for", authConfig.Email)

	return nil
}

func validateLogin(email string, pass string) (string, error) {
	signin, err := httpclientwrapper.NewSignin("https", utils2.GetHost())
	if err != nil {
		return "", fmt.Errorf("unable to login, %s", err.Error())
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
			return "", fmt.Errorf("unable to login, %s", err.Error())
		}

		if serviceErr.Name == "auth_failure" || serviceErr.Name == "bad_request" {
			return "", fmt.Errorf("unable to login, either email or password is incorrect")
		}

		return "", fmt.Errorf("unable to login, %s", serviceErr.Name)
	}
	return res.JWTToken, nil
}

func resetLogin() {
	email = ""
	password = ""
	passwordStdin = false
}
