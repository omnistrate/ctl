package login

import (
	"fmt"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	ctlutils "github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"strings"
)

func PasswordLogin(cmd *cobra.Command, args []string) error {
	if len(password) > 0 {
		ctlutils.PrintWarning("Notice: Using the --password flag is insecure. Please consider using the --password-stdin flag instead. Refer to the help documentation for examples.")
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

	token, err := dataaccess.LoginWithPassword(email, password)
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

	ctlutils.PrintSuccess("Successfully logged in")

	return nil
}
