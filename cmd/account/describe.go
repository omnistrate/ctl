package account

import (
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	describeExample = `  # Describe account with name
  omctl account describe [account-name]

  # Describe account with ID
  omctl account describe --id=[account-id]`
)

var describeCmd = &cobra.Command{
	Use:     "describe [account-name] [flags]",
	Short:   "Describe an account",
	Long:    "This command helps you describe an account.",
	Example: describeExample,
	RunE:    runDescribe,
	PostRun: func(cmd *cobra.Command, args []string) {
		dataaccess.AskVerifyAccountIfAny()
	},
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.MaximumNArgs(1) // Require at most 1 argument

	describeCmd.Flags().String("id", "", "Account ID")
	describeCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported.") // Override inherited flag
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

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
	err = validateDescribeArguments(args, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := utils.GetToken()
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
	id, name, err = getAccount(token, name, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe account
	account, err := dataaccess.DescribeAccount(token, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully retrieved account details")

	// Print output
	err = utils.PrintTextTableJsonOutput(output, account)

	return nil
}

// Helper functions

func validateDescribeArguments(args []string, accountIDArg string) error {
	if len(args) == 0 && accountIDArg == "" {
		return errors.New("account name or ID must be provided")
	}

	if len(args) != 0 && accountIDArg != "" {
		return errors.New("only one of account name or ID can be provided")
	}

	return nil
}

func getAccount(token, accountNameArg, accountIDArg string) (accountID, accountName string, err error) {
	// List accounts
	listRes, err := dataaccess.ListAccounts(token, "all")
	if err != nil {
		return
	}

	count := 0
	if accountNameArg != "" {
		for _, account := range listRes.AccountConfigs {
			if strings.EqualFold(account.Name, accountNameArg) {
				accountID = string(account.ID)
				accountName = account.Name
				count++
			}
		}
	} else {
		for _, account := range listRes.AccountConfigs {
			if strings.EqualFold(string(account.ID), accountIDArg) {
				accountID = string(account.ID)
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
