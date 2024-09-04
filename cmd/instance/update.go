package instance

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	updateExample = `# Update an instance deployment
omctl instance update instance-abcd1234`
)

var updateCmd = &cobra.Command{
	Use:          "update [instance-id]",
	Short:        "Update an instance deployment for your service",
	Long:         `This command helps you update the instance for your service.`,
	Example:      updateExample,
	RunE:         runUpdate,
	SilenceUsage: true,
}

func init() {

	updateCmd.Flags().String("param", "", "Parameters for the instance deployment")
	updateCmd.Flags().String("param-file", "", "Json file containing parameters for the instance deployment")

	if err := updateCmd.MarkFlagFilename("param-file"); err != nil {
		return
	}

	updateCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func runUpdate(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	instanceID := args[0]

	// Retrieve flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	param, err := cmd.Flags().GetString("param")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	paramFile, err := cmd.Flags().GetString("param-file")
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
		msg := "Updating instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, resourceID, err := getInstance(token, instanceID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Format parameters
	formattedParams, err := formatParams(param, paramFile)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Update instance
	err = dataaccess.UpdateResourceInstance(token, inventoryapi.FleetUpdateResourceInstanceRequest{
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
		ResourceID:    inventoryapi.ResourceID(resourceID),
		RequestParams: formattedParams,
	})
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully updated instance")

	// Search for the instance
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", instanceID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if len(searchRes.ResourceInstanceResults) == 0 {
		err = errors.New("failed to find the updated instance")
		utils.PrintError(err)
		return err
	}

	// Format instance
	formattedInstance := formatInstance(searchRes.ResourceInstanceResults[0], false)

	// Print output
	if err = utils.PrintTextTableJsonOutput(output, formattedInstance); err != nil {
		return err
	}

	return nil
}
