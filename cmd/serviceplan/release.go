package serviceplan

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	releaseExample = `  # Release service plan
  omctl service-plan release [service-name] [plan-name]

  # Release service plan by ID instead of name
  omctl service-plan release --service-id [service-id] --plan-id [plan-id]`
)

var releaseCmd = &cobra.Command{
	Use:          "release [service-name] [plan-name] [flags]",
	Short:        "Release a service plan",
	Long:         `This command helps you release a service plan from your service.`,
	Example:      releaseExample,
	RunE:         runRelease,
	SilenceUsage: true,
}

func init() {
	releaseCmd.Flags().String("release-description", "", "Set custom release description for this release version")
	releaseCmd.Flags().Bool("release-as-preferred", false, "Release the service plan as preferred")
	releaseCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	releaseCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	releaseCmd.Flags().StringP("plan-id", "", "", "Plan ID. Required if plan name is not provided")
}

func runRelease(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	releaseDescription, _ := cmd.Flags().GetString("release-description")
	releaseAsPreferred, _ := cmd.Flags().GetBool("release-as-preferred")
	output, _ := cmd.Flags().GetString("output")
	serviceID, _ := cmd.Flags().GetString("service-id")
	planID, _ := cmd.Flags().GetString("plan-id")

	// Validate input arguments
	if err := validateReleaseArguments(args, serviceID, planID); err != nil {
		utils.PrintError(err)
		return err
	}

	// Set service and plan names if provided in args
	var serviceName, planName string
	if len(args) == 2 {
		serviceName, planName = args[0], args[1]
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
		msg := "Releasing service plan..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if service plan exists
	serviceID, _, planID, _, _, err = getServicePlan(token, serviceID, serviceName, planID, planName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get service api id
	productTier, err := dataaccess.DescribeProductTier(token, serviceID, planID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	serviceModel, err := dataaccess.DescribeServiceModel(token, serviceID, string(productTier.ServiceModelID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	serviceAPIID := string(serviceModel.ServiceAPIID)

	// Release service plan
	err = dataaccess.ReleaseServicePlan(token, serviceID, serviceAPIID, planID, getReleaseDescription(releaseDescription), releaseAsPreferred)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully released service plan")

	// Get the service plan details
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("serviceplan:%s", planID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	targetVersion, err := dataaccess.FindLatestVersion(token, serviceID, planID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var targetServicePlan *inventoryapi.ServicePlanSearchRecord
	for _, servicePlan := range searchRes.ServicePlanResults {
		if string(servicePlan.ServiceID) != serviceID || servicePlan.ID != planID || servicePlan.Version != targetVersion {
			continue
		}
		targetServicePlan = servicePlan
	}

	// Format output
	formattedServicePlan, err := formatServicePlanVersion(targetServicePlan, false)
	if err != nil {
		return err
	}

	// Marshal data
	data, err := json.MarshalIndent(formattedServicePlan, "", "    ")
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

func validateReleaseArguments(args []string, serviceID, planID string) error {
	if len(args) == 0 && (serviceID == "" || planID == "") {
		return fmt.Errorf("please provide the service name and service plan name or the service ID and service plan ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [plan-name]", strings.Join(args, " "))
	}
	return nil
}

func getReleaseDescription(releaseDescription string) *string {
	if releaseDescription == "" {
		return nil
	}
	return &releaseDescription
}
