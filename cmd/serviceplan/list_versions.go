package serviceplan

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listVersionsExample = `# List service plan versions of the service postgres in the prod and dev environments
omnistrate service-plan list-versions postgres postgres -o=table -f="service_name:postgres,environment:prod" -f="service:postgres,environment:dev"`
)

var listVersionsCmd = &cobra.Command{
	Use:   "list-versions [service-name] [plan-name] [flags]",
	Short: "List service plan versions for your services",
	Long: `This command helps you list service plan versions for your services.
You can filter for specific service plan versions by using the filter flag.`,
	Example:      listVersionsExample,
	RunE:         runListVersions,
	SilenceUsage: true,
}

func init() {
	listVersionsCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	listVersionsCmd.Flags().StringP("plan-id", "", "", "Environment ID. Required if plan name is not provided")
	listVersionsCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listVersionsCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of service plan versions. E.g.: key1:value1,key2:value2, which filters service plans where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.ServicePlanVersion{}), ",")+". Check the examples for more details.")
	listVersionsCmd.Flags().Bool("truncate", false, "Truncate long names in the output")
}

func runListVersions(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve command-line flags
	serviceId, _ := cmd.Flags().GetString("service-id")
	planId, _ := cmd.Flags().GetString("plan-id")
	output, _ := cmd.Flags().GetString("output")
	filters, _ := cmd.Flags().GetStringArray("filter")
	truncateNames, _ := cmd.Flags().GetBool("truncate")

	// Validate input arguments
	if err := validateListVersionsArguments(args, serviceId, planId); err != nil {
		utils.PrintError(err)
		return err
	}

	// Parse and validate filters
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.ServicePlanVersion{}))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Set service and service plan names if provided in args
	var serviceName, planName string
	if len(args) == 2 {
		serviceName, planName = args[0], args[1]
	}

	// Ensure user is logged in
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
		spinner = sm.AddSpinner("Listing service plan versions...")
		sm.Start()
	}

	// Check if the service plan exists
	_, _, planId, _, _, err = getServicePlan(token, serviceId, serviceName, planId, planName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Search service plans versions
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("serviceplan:%s", planId))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var servicePlanVersions []string

	// Process and filter service plans
	for _, servicePlanVersion := range searchRes.ServicePlanResults {
		formattedServicePlanVersion, err := formatServicePlanVersion(servicePlanVersion, truncateNames)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		match, err := utils.MatchesFilters(formattedServicePlanVersion, filterMaps)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		data, err := json.MarshalIndent(formattedServicePlanVersion, "", "    ")
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		if match {
			servicePlanVersions = append(servicePlanVersions, string(data))
		}
	}

	// Handle case when no service plans match
	if len(servicePlanVersions) == 0 {
		utils.HandleSpinnerSuccess(spinner, sm, "No service plan versions found.")
		return nil
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Service plan versions retrieved successfully")

	// Format output as requested
	err = utils.PrintTextTableJsonArrayOutput(output, servicePlanVersions)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func formatServicePlanVersion(servicePlan *inventoryapi.ServicePlanSearchRecord, truncateNames bool) (model.ServicePlanVersion, error) {
	serviceName := servicePlan.ServiceName
	envName := servicePlan.ServiceEnvironmentName
	planName := servicePlan.Name

	if truncateNames {
		serviceName = utils.TruncateString(serviceName, defaultMaxNameLength)
		envName = utils.TruncateString(envName, defaultMaxNameLength)
		planName = utils.TruncateString(planName, defaultMaxNameLength)
	}

	var releaseDescription string
	if servicePlan.VersionName != nil {
		releaseDescription = *servicePlan.VersionName
	}

	return model.ServicePlanVersion{
		PlanID:             servicePlan.ID,
		PlanName:           planName,
		ServiceID:          string(servicePlan.ServiceID),
		ServiceName:        serviceName,
		Environment:        envName,
		Version:            servicePlan.Version,
		ReleaseDescription: releaseDescription,
		VersionSetStatus:   servicePlan.VersionSetStatus,
	}, nil
}

func validateListVersionsArguments(args []string, serviceId, planId string) error {
	if len(args) == 0 && (serviceId == "" || planId == "") {
		return fmt.Errorf("please provide the service name and service plan name or the service ID and service plan ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [plan-name]", strings.Join(args, " "))
	}
	return nil
}
