package account

import (
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	deleteExample = `# Delete account with name
omctl account delete [account-name]

# Delete account with ID
omctl account delete --id=[account-ID]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [account-name] [flags]",
	Short:        "Delete a Cloud Provider Account",
	Long:         `This command helps you delete a Cloud Provider Account from your account list.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Args = cobra.MaximumNArgs(1) // Require at most one argument

	deleteCmd.Flags().String("id", "", "Account ID")
}

func runDelete(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	var name string
	if len(args) > 0 {
		name = args[0]
	}

	// Retrieve flags
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate input args
	err = validateDeleteArguments(args, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Deleting account..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if account exists
	id, _, err = getAccount(token, name, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Delete account
	err = dataaccess.DeleteAccount(token, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully deleted account")

	return nil
}

// Helper functions

func validateDeleteArguments(args []string, accountIDArg string) error {
	if len(args) == 0 && accountIDArg == "" {
		return errors.New("account name or ID must be provided")
	}

	if len(args) != 0 && accountIDArg != "" {
		return errors.New("only one of account name or ID can be provided")
	}

	return nil
}
