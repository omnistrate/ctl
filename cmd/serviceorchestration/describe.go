package serviceorchestration

import (
	"errors"

	"github.com/omnistrate/ctl/cmd/common"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe an service orchestration deployment
omctl service-orchestration describe so-abcd1234`
)

var describeCmd = &cobra.Command{
	Use:          "describe [so-id]",
	Short:        "Describe an service orchestration deployment",
	Long:         `This command helps you describe a service orchestration deployment.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	describeCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	soID := args[0]

	// Retrieve flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate output flag
	if output != "json" {
		err = errors.New("only json output is supported")
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Describing instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Describe service orchestration
	serviceOrchestration, err := dataaccess.DescribeServicesOrchestration(
		cmd.Context(),
		token,
		soID,
	)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully described service orchestration deployment")
	// Print output
	err = utils.PrintTextTableJsonOutput(output, serviceOrchestration)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
