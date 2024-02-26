package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/omnistrate/api-design/v1/pkg/registration/gen/http/signin_api/client"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	"github.com/spf13/cobra"
	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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
	Example: `  cat ~/omnistrate_pass.txt | omnistrate-cli login -u user --password-stdin
	  echo $PASSWORD | omnistrate-cli auth login -s 
	  omnistrate-cli login -u user -p password`,
	RunE: runLogin,
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().StringVarP(&email, "email", "e", "", "email")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	loginCmd.Flags().BoolVarP(&passwordStdin, "password-stdin", "s", false, "Reads the password from stdin")
	loginCmd.Flags().Duration("timeout", 5*time.Second, "Override the timeout for this API call")
}

func runLogin(cmd *cobra.Command, args []string) error {

	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		return err
	}

	if len(email) == 0 {
		return fmt.Errorf("must provide --email or -u")
	}

	if len(password) > 0 {
		fmt.Println("WARNING! Using --password is insecure, consider using: cat ~/omnistrate_pass.txt | omnistrate-cli login -u user --password-stdin")
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

		var passwordFromStdin []byte
		if passwordFromStdin, err = io.ReadAll(os.Stdin); err != nil {
			return err
		}
		password = strings.TrimSpace(string(passwordFromStdin))
	}

	password = strings.TrimSpace(password)
	if len(password) == 0 {
		return fmt.Errorf("must provide a non-empty password via --password or --password-stdin")
	}

	fmt.Println("Calling the Omnistrate server to validate the credentials...")

	var token string
	if token, err = validateLogin(email, password, timeout); err != nil {
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

	fmt.Println("credentials saved for", authConfig.Email)

	return nil
}

func validateLogin(email string, pass string, timeout time.Duration) (string, error) {
	cli := client.NewClient("https", "api.omnistrate.cloud", &http.Client{Timeout: timeout}, goahttp.RequestEncoder, goahttp.ResponseDecoder, false)
	request := signinapi.SigninRequest{
		Email:          email,
		HashedPassword: utils.HashPassword(pass),
	}
	endpoint := cli.Signin()
	res, err := endpoint(context.Background(), &request)

	if err == nil {
		result := res.(*signinapi.SigninResult)
		return result.JWTToken, nil
	}

	var serviceErr *goa.ServiceError
	ok := errors.As(err, &serviceErr)
	if !ok {
		return "", fmt.Errorf("unable to login, %s", err.Error())
	}

	if serviceErr.Name == "forbidden" {
		return "", fmt.Errorf("unable to login, either username or password is incorrect")
	}

	return "", fmt.Errorf("unable to login, %s", err.Error())
}
