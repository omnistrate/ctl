package domain

import (
	"encoding/json"
	saasportalapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/saas_portal_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listExample = `  # List domains
  omctl domain list -o=table`
)

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List SaaS Portal custom domains",
	Long: `This command helps you list SaaS Portal custom domains.
You can filter for specific domains by using the filter flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of domains. E.g.: key1:value1,key2:value2, which filters domains where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Domain{}), ",")+". Check the examples for more details.")
}

func runList(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve command-line flags
	output, _ := cmd.Flags().GetString("output")
	filters, _ := cmd.Flags().GetStringArray("filter")

	// Parse and validate filters
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.Domain{}))
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

	// Retrieve domains and domains
	listRes, err := dataaccess.ListDomains(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var domainsData []string

	// Process and filter domains
	for _, domain := range listRes.CustomDomains {
		formattedDomain, err := formatDomain(domain)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		match, err := utils.MatchesFilters(formattedDomain, filterMaps)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		data, err := json.MarshalIndent(formattedDomain, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if match {
			domainsData = append(domainsData, string(data))
		}
	}

	// Handle case when no domains match
	if len(domainsData) == 0 {
		utils.PrintInfo("No domains found.")
		return nil
	}

	// Format output as requested
	err = utils.PrintTextTableJsonArrayOutput(output, domainsData)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func formatDomain(domain *saasportalapi.CustomDomain) (model.Domain, error) {
	return model.Domain{
		EnvironmentType: string(domain.EnvironmentType),
		Name:            domain.Name,
		Domain:          domain.CustomDomain,
		Status:          string(domain.Status),
		ClusterEndpoint: domain.ClusterEndpoint,
	}, nil
}
