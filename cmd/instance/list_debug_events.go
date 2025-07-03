package instance

import (
	"errors"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	listDebugEventsExample = `# List debug events for an instance
omctl instance list-debug-events i-1234`
)

var listDebugEventsCmd = &cobra.Command{
	Use:          "list-debug-events [instance-id]",
	Short:        "List debug events for an instance deployment",
	Long:         `This command helps you list debug events for an instance deployment that has debug mode enabled.`,
	Example:      listDebugEventsExample,
	RunE:         runListDebugEvents,
	SilenceUsage: true,
}

func init() {
	listDebugEventsCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	listDebugEventsCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")
}

func runListDebugEvents(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]

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
		msg := "Listing debug events for instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// List debug events for the instance
	events, err := dataaccess.ListInstanceEvents(cmd.Context(), token, serviceID, environmentID, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully retrieved debug events for instance")

	// Print output
	err = utils.PrintTextTableJsonOutput(output, events)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
