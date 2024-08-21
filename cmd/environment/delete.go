package environment

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	deleteExample = `# Delete environment
omnistrate environment delete [service-name] [environment-name]

# Delete environment by ID instead of name
omnistrate environment delete --service-id [service-id] --environment-id [environment-id]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [service-name] [environment-name] [flags]",
	Short:        "Delete a environment",
	Long:         `This command helps you delete a environment in your service.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	deleteCmd.Flags().StringP("environment-id", "", "", "Environment ID. Required if environment name is not provided")
}
func runDelete(cmd *cobra.Command, args []string) error {
	defer cleanUpDeleteFlagsAndArgs(cmd, &args)

	// Retrieve flags
	serviceId, _ := cmd.Flags().GetString("service-id")
	environmentId, _ := cmd.Flags().GetString("environment-id")

	// Validate input arguments
	if err := validateDeleteArguments(args, serviceId, environmentId); err != nil {
		utils.PrintError(err)
		return err
	}

	// Set service and environment names if provided in args
	var serviceName, environmentName string
	if len(args) == 2 {
		serviceName, environmentName = args[0], args[1]
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
	if cmd.Flag("output").Value.String() != "json" {
		sm = ysmrr.NewSpinnerManager()
		spinner = sm.AddSpinner("Deleting environment...")
		sm.Start()
	}

	// Check if the environment exists
	services, err := dataaccess.ListServices(token)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	serviceId, serviceName, environmentId, environmentName, err = getServiceEnvironment(services, serviceId, serviceName, environmentId, environmentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Delete the environment
	if err = dataaccess.DeleteServiceEnvironment(token, serviceId, environmentId); err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Handle success message
	if spinner != nil {
		spinner.UpdateMessage("Successfully deleted environment")
		spinner.Complete()
		sm.Stop()
	}

	return nil
}

// Helper functions

func validateDeleteArguments(args []string, serviceId, environmentId string) error {
	if len(args) == 0 && (serviceId == "" || environmentId == "") {
		return fmt.Errorf("please provide the service name and environment name or the service ID and environment ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [environment-name]", strings.Join(args, " "))
	}
	return nil
}

func cleanUpDeleteFlagsAndArgs(cmd *cobra.Command, args *[]string) {
	// Clean up flags
	_ = cmd.Flags().Set("service-id", "")
	_ = cmd.Flags().Set("environment-id", "")

	// Clean up arguments by resetting the slice to nil or an empty slice
	*args = nil
}
