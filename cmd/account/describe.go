package account

import (
	"context"
	"strings"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe account with name
omctl account describe [account-name]

# Describe account with ID
omctl account describe --id=[account-id]`
)

var describeCmd = &cobra.Command{
	Use:          "describe [account-name] [flags]",
	Short:        "Describe a Cloud Provider Account",
	Long:         "This command helps you get details of a cloud provider account.",
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.MaximumNArgs(1) // Require at most 1 argument

	describeCmd.Flags().String("id", "", "Account ID")
	describeCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported.") // Override inherited flag
}

func runDescribe(cmd *cobra.Command, args []string) error {
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
	err = validateDescribeArguments(args, id, output)
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
	id, _, err = getAccount(cmd.Context(), token, name, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe account
	account, err := dataaccess.DescribeAccount(cmd.Context(), token, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully retrieved account details")

	// Print output
	err = utils.PrintTextTableJsonOutput(output, account)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Ask user to verify account if output is not JSON
	if output != "json" {
		dataaccess.AskVerifyAccountIfAny(cmd.Context())
	}

	return nil
}

// Helper functions

func validateDescribeArguments(args []string, accountIDArg, output string) error {
	if len(args) == 0 && accountIDArg == "" {
		return errors.New("account name or ID must be provided")
	}

	if len(args) != 0 && accountIDArg != "" {
		return errors.New("only one of account name or ID can be provided")
	}

	if output != "json" {
		return errors.New("only json output is supported")
	}

	return nil
}

func getAccount(ctx context.Context, token, accountNameArg, accountIDArg string) (accountID, accountName string, err error) {
	// List accounts
	listRes, err := dataaccess.ListAccounts(ctx, token, "all")
	if err != nil {
		return
	}

	count := 0
	if accountNameArg != "" {
		for _, account := range listRes.AccountConfigs {
			if strings.EqualFold(account.Name, accountNameArg) {
				accountID = account.Id
				accountName = account.Name
				count++
			}
		}
	} else {
		for _, account := range listRes.AccountConfigs {
			if strings.EqualFold(account.Id, accountIDArg) {
				accountID = account.Id
				accountName = account.Name
				count++
			}
		}
	}

	if count == 0 {
		err = errors.New("account not found")
		return
	}

	if count > 1 {
		err = errors.New("multiple accounts found with the same name. Please specify account ID")
		return
	}

	return
}
