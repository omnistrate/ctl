package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/omnistrate-oss/ctl/cmd/common"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate-oss/ctl/internal/config"
	"github.com/omnistrate-oss/ctl/internal/dataaccess"
	"github.com/omnistrate-oss/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe an instance deployment
omctl instance describe instance-abcd1234`
)

type InstanceStatusType string

var InstanceStatus InstanceStatusType

const (
	InstanceStatusRunning   InstanceStatusType = "RUNNING"
	InstanceStatusStopped   InstanceStatusType = "STOPPED"
	InstanceStatusFailed    InstanceStatusType = "FAILED"
	InstanceStatusCancelled InstanceStatusType = "CANCELLED"
	InstanceStatusUnknown   InstanceStatusType = "UNKNOWN"
)

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
	var instance *openapiclientfleet.ResourceInstance
	instance, err = dataaccess.DescribeResourceInstance(cmd.Context(), token, serviceID, environmentID, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully described instance")
	if instance.ConsumptionResourceInstanceResult.Status != nil {
		InstanceStatus = InstanceStatusType(*instance.ConsumptionResourceInstanceResult.Status)
	} else {
		InstanceStatus = InstanceStatusUnknown
	}

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

func getResourceFromInstance(ctx context.Context, token string, instanceID string, resourceName string) (resourceID, resourceType string, err error) {
	// Check if instance exists
	serviceID, environmentID, _, _, err := getInstance(ctx, token, instanceID)
	if err != nil {
		return
	}

	// Retrieve resource ID
	instanceDes, err := dataaccess.DescribeResourceInstance(ctx, token, serviceID, environmentID, instanceID)
	if err != nil {
		return
	}

	versionSetDes, err := dataaccess.DescribeVersionSet(ctx, token, serviceID, instanceDes.ProductTierId, instanceDes.TierVersion)
	if err != nil {
		return
	}

	for _, resource := range versionSetDes.Resources {
		if resource.Name == resourceName {
			resourceID = resource.Id
			if resource.ManagedResourceType != nil {
				resourceType = strings.ToLower(*resource.ManagedResourceType)
			}
		}
	}

	return
}
