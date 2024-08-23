package service

import (
	"encoding/json"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listExample = `# List services
omnistrate service list -o=table`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List services for your account",
	Long: `This command helps you list services for your account.
You can filter for specific services by using the filters flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of services. E.g.: key1:value1,key2:value2, which filters services where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Service{}), ",")+". Check the examples for more details.")
	listCmd.Flags().Bool("truncate", false, "Truncate long names in the output")
}

func runList(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve command-line flags
	output, _ := cmd.Flags().GetString("output")
	filters, _ := cmd.Flags().GetStringArray("filter")
	truncateNames, _ := cmd.Flags().GetBool("truncate")

	// Parse and validate filters
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.Service{}))
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

	// Retrieve services and services
	listRes, err := dataaccess.ListServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var servicesData []string

	// Process and filter services
	for _, service := range listRes.Services {
		formattedService, err := formatService(service, truncateNames)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		match, err := utils.MatchesFilters(formattedService, filterMaps)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		data, err := json.MarshalIndent(formattedService, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if match {
			servicesData = append(servicesData, string(data))
		}
	}

	// Handle case when no services match
	if len(servicesData) == 0 {
		utils.PrintInfo("No services found.")
		return nil
	}

	// Format output as requested
	err = utils.PrintTextTableJsonArrayOutput(output, servicesData)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func formatService(service *serviceapi.DescribeServiceResult, truncateNames bool) (model.Service, error) {
	// Retrieve environments
	environments := make([]string, 0)
	for _, environment := range service.ServiceEnvironments {
		if environment == nil {
			continue
		}
		environments = append(environments, environment.Name)
	}

	// Truncate service name if requested
	serviceName := service.Name
	if truncateNames {
		serviceName = utils.TruncateString(serviceName, defaultMaxNameLength)
	}

	return model.Service{
		ID:           string(service.ID),
		Name:         serviceName,
		Environments: strings.Join(environments, ","),
	}, nil
}
