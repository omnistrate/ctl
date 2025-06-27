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
	getDebugEventExample = `# Get details of a specific debug event for an instance
omctl instance get-debug-event i-1234 --event-id event-5678`
)

var getDebugEventCmd = &cobra.Command{
	Use:          "get-debug-event [instance-id] --event-id [event-id]",
	Short:        "Get details of a specific debug event for an instance deployment",
	Long:         `This command helps you get detailed information about a specific debug event for an instance deployment.`,
	Example:      getDebugEventExample,
	RunE:         runGetDebugEvent,
	SilenceUsage: true,
}

func init() {
	getDebugEventCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	getDebugEventCmd.Flags().StringP("event-id", "e", "", "Event ID")
	getDebugEventCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	var err error
	if err = getDebugEventCmd.MarkFlagRequired("event-id"); err != nil {
		return
	}
}

func runGetDebugEvent(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]

	// Retrieve flags
	eventID, err := cmd.Flags().GetString("event-id")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if eventID == "" {
		err = errors.New("event id is required")
		utils.PrintError(err)
		return err
	}

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
		msg := "Getting debug event details..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get debug event details
	event, err := dataaccess.DescribeInstanceEvent(cmd.Context(), token, serviceID, environmentID, instanceID, eventID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully retrieved debug event details")

	// Print output
	err = utils.PrintTextTableJsonOutput(output, event)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}