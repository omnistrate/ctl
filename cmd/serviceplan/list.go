package serviceplan

import (
	"encoding/json"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listExample = `# List service plans of the service postgres in the prod and dev environments
omnistrate service-plan list -o=table -f="service_name:postgres,environment:prod" -f="service:postgres,environment:dev"`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List service plans for your services",
	Long: `This command helps you list service plans for your services.
You can filter for specific service plans by using the filters flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of service plans. E.g.: key1:value1,key2:value2, which filters service plans where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Environment{}), ",")+". Check the examples for more details.")
	listCmd.Flags().Bool("truncate", false, "Truncate long names in the output")
}

func runList(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve command-line flags
	output, _ := cmd.Flags().GetString("output")
	filters, _ := cmd.Flags().GetStringArray("filter")
	truncateNames, _ := cmd.Flags().GetBool("truncate")

	// Parse and validate filters
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.ServicePlan{}))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Ensure user is logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Search service plans
	searchRes, err := dataaccess.SearchInventory(token, "serviceplan:pt")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var servicePlans []string

	// Process and filter service plans
	for _, servicePlan := range searchRes.ServicePlanResults {
		env, err := formatServicePlan(servicePlan, truncateNames)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		match, err := utils.MatchesFilters(env, filterMaps)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		data, err := json.MarshalIndent(env, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if match {
			servicePlans = append(servicePlans, string(data))
		}
	}

	// Handle case when no service plans match
	if len(servicePlans) == 0 {
		utils.PrintInfo("No service plans found.")
		return nil
	}

	// Format output as requested
	err = utils.PrintTextTableJsonArrayOutput(output, servicePlans)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func formatServicePlan(servicePlan *inventoryapi.ServicePlanSearchRecord, truncateNames bool) (model.ServicePlan, error) {
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

	return model.ServicePlan{
		PlanID:             servicePlan.ID,
		PlanName:           planName,
		ServiceID:          string(servicePlan.ServiceID),
		ServiceName:        serviceName,
		Environment:        envName,
		Version:            servicePlan.Version,
		ReleaseDescription: releaseDescription,
		VersionSetStatus:   servicePlan.VersionSetStatus,
		DeploymentType:     string(servicePlan.DeploymentType),
		TenancyType:        string(servicePlan.TenancyType),
	}, nil
}
