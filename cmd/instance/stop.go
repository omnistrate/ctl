package instance

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	stopExample = `  # Stop an instance deployment
  omctl instance stop instance-abcd1234`
)

var stopCmd = &cobra.Command{
	Use:          "stop [instance-id]",
	Short:        "Stop an instance deployment for your service",
	Long:         `This command helps you stop the instance for your service.`,
	Example:      stopExample,
	RunE:         runStop,
	SilenceUsage: true,
}

func init() {
	stopCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	stopCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func runStop(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	instanceID := args[0]

	// Retrieve flags
	output, err := cmd.Flags().GetString("output")
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
		msg := "Stoping instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, resourceID, err := getInstance(token, instanceID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Stop instance
	err = dataaccess.StopResourceInstance(token, serviceID, environmentID, resourceID, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully stopped instance")

	// Search for the instance
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", instanceID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if len(searchRes.ResourceInstanceResults) == 0 {
		err = errors.New("failed to find the stopped instance")
		utils.PrintError(err)
		return err
	}

	// Format instance
	formattedInstance := formatInstance(searchRes.ResourceInstanceResults[0], false)

	// Marshal instance to JSON
	data, err := json.MarshalIndent(formattedInstance, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	if err = utils.PrintTextTableJsonOutput(output, string(data)); err != nil {
		return err
	}

	return nil
}
