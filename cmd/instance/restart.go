package instance

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	restartExample = `  # Restart an instance deployment
  omctl instance restart instance-abcd1234`

	defaultRestartOutput = "json"
)

var restartCmd = &cobra.Command{
	Use:          "restart [instance-id]",
	Short:        "Restart an instance deployment for your service",
	Long:         `This command helps you restart the instance for your service.`,
	Example:      restartExample,
	RunE:         runRestart,
	SilenceUsage: true,
}

func init() {
	restartCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func runRestart(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	instanceID := args[0]
	output := defaultRestartOutput

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
		msg := "Restarting instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, resourceID, err := getInstance(token, instanceID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Restart instance
	err = dataaccess.RestartResourceInstance(token, serviceID, environmentID, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully created instance")

	// Marshal instance to JSON
	data, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	fmt.Println(string(data))

	return nil
}
