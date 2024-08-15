package instance

import (
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	listExample = `# List instances of the service postgres in the prod and dev environments
omnistrate instance list -o=table -f="service:postgres,environment:PROD" -f="service:postgres,environment:DEV"`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List instance deployments for your services",
	Long: `This command helps you list instance deployments for your services.
You can filter for specific instances by using the filters flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of instances. E.g.: key1:value1,key2:value2, which filters instances where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Instance{}), ",")+". Check the examples for more details.")
	listCmd.Flags().Bool("truncate", false, "Truncate long names in the output")
}

func runList(cmd *cobra.Command, args []string) error {
	// Get flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	filters, err := cmd.Flags().GetStringArray("filter")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	truncateNames, err := cmd.Flags().GetBool("truncate")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Parse filters into a map
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.Instance{}))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Get all instances
	searchRes, err := dataaccess.SearchInventory(token, "resourceinstance:i")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	instances := make([]model.Instance, 0)
	for _, instance := range searchRes.ResourceInstanceResults {
		if instance == nil {
			continue
		}
		planName := ""
		if instance.ProductTierName != nil {
			planName = *instance.ProductTierName
		}
		planVersion := ""
		if instance.ProductTierVersion != nil {
			planVersion = *instance.ProductTierVersion
		}
		envType := ""
		if instance.ServiceEnvironmentType != nil {
			envType = string(*instance.ServiceEnvironmentType)
		}
		serviceName := instance.ServiceName
		if truncateNames {
			serviceName = utils.TruncateString(serviceName, defaultMaxNameLength)
			planName = utils.TruncateString(planName, defaultMaxNameLength)
		}
		formattedInstance := model.Instance{
			ID:            instance.ID,
			Service:       serviceName,
			Environment:   envType,
			Plan:          planName,
			Version:       planVersion,
			Resource:      instance.ResourceName,
			CloudProvider: string(instance.CloudProvider),
			Region:        instance.RegionCode,
			Status:        string(instance.Status),
		}

		// Check if the instance matches the filters
		ok, err := utils.MatchesFilters(formattedInstance, filterMaps)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		if ok {
			instances = append(instances, formattedInstance)
		}
	}

	var jsonData []string
	for _, instance := range instances {
		data, err := json.MarshalIndent(instance, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		jsonData = append(jsonData, string(data))
	}

	if len(jsonData) == 0 {
		utils.PrintInfo("No instances found.")
		return nil
	}

	switch output {
	case "text":
		err = utils.PrintText(jsonData)
		if err != nil {
			return err
		}
	case "table":
		err = utils.PrintTable(jsonData)
		if err != nil {
			return err
		}
	case "json":
		fmt.Printf("%+v\n", jsonData)
	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
