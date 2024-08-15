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
	listExample = `# List instances of the service postgres in the prod environment
omnistrate instance list --output=table --filters="service:postgres,environment:prod"`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var supportedFilterKeys = []string{"id", "service", "environment", "plan", "version", "resource", "cloud_provider", "region", "status"}

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
	listCmd.Flags().StringArrayP("filters", "f", []string{}, "Filters to apply to the list of instances. Format: 'key:value,key:value' 'key:value'. Supported keys: "+strings.Join(supportedFilterKeys, ","))
	listCmd.Flags().Bool("truncate", false, "Truncate long names in the output")
}

func runList(cmd *cobra.Command, args []string) error {
	// Get flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	filters, err := cmd.Flags().GetStringArray("filters")
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
	filterMaps, err := utils.ParseFilters(filters, supportedFilterKeys)
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

	// Check if the instance exists
	searchRes, err := dataaccess.SearchInventory(token, "resourceinstance:i") // Get all instances
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
		var tableWriter *utils.Table
		if tableWriter, err = utils.NewTableFromJSONTemplate(json.RawMessage(jsonData[0])); err != nil {
			// Just print the JSON directly and return
			fmt.Printf("%+v\n", jsonData)
			return err
		}

		for _, data := range jsonData {
			if err = tableWriter.AddRowFromJSON(json.RawMessage(data)); err != nil {
				// Just print the JSON directly and return
				fmt.Printf("%+v\n", jsonData)
				return err
			}
		}

		tableWriter.Print()
	case "json":
		fmt.Printf("%+v\n", jsonData)
	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
