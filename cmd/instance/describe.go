package instance

import (
	"context"
	"errors"
	"fmt"

	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe an instance deployment
omctl instance describe instance-abcd1234`
)

var InstanceStatus string

var describeCmd = &cobra.Command{
	Use:          "describe [instance-id]",
	Short:        "Describe an instance deployment for your service",
	Long:         `This command helps you describe the instance for your service.`,
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
		msg := "Describing instance..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe instance
	var instance *inventoryapi.ResourceInstance
	instance, err = dataaccess.DescribeInstance(cmd.Context(), token, serviceID, environmentID, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully described instance")
	InstanceStatus = string(instance.ConsumptionResourceInstanceResult.Status)

	// Print output
	err = utils.PrintTextTableJsonOutput(output, instance)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}

// Helper functions

func getInstance(ctx context.Context, token, instanceID string) (serviceID, environmentID, productTierID, resourceID string, err error) {
	searchRes, err := dataaccess.SearchInventory(ctx, token, fmt.Sprintf("resourceinstance:%s", instanceID))
	if err != nil {
		return
	}

	var found bool
	for _, instance := range searchRes.ResourceInstanceResults {
		if instance.Id == instanceID {
			serviceID = instance.ServiceId
			environmentID = instance.ServiceEnvironmentId
			productTierID = instance.ProductTierId
			if instance.ResourceId != nil {
				resourceID = *instance.ResourceId
			}
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID)
		return
	}

	return
}
