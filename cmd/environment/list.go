package environment

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
	listExample = `# List environments of the service postgres in the prod and dev environment types
omnistrate environment list -o=table -f="service_name:postgres,environment_type:PROD" -f="service:postgres,environment_type:DEV"`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments for your services",
	Long: `This command helps you list environments for your services.
You can filter for specific environments by using the filters flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of environments. E.g.: key1:value1,key2:value2, which filters environments where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Environment{}), ",")+". Check the examples for more details.")
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
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.Environment{}))
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

	// Get all environments
	services, err := dataaccess.ListServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	environments := make([]model.Environment, 0)
	for _, service := range services.Services {
		for _, environment := range service.ServiceEnvironments {
			if environment == nil {
				continue
			}
			serviceName := service.Name
			if truncateNames {
				serviceName = utils.TruncateString(serviceName, defaultMaxNameLength)
			}
			envName := environment.Name
			if truncateNames {
				envName = utils.TruncateString(envName, defaultMaxNameLength)
			}
			envType := ""
			if environment.Type != nil {
				envType = string(*environment.Type)
			}
			sourceEnvName := ""
			if environment.SourceEnvironmentName != nil {
				sourceEnvName = *environment.SourceEnvironmentName
			}
			formattedEnvironment := model.Environment{
				EnvironmentID:   string(environment.ID),
				EnvironmentName: envName,
				EnvironmentType: envType,
				ServiceID:       string(service.ID),
				ServiceName:     serviceName,
				SourceEnvName:   sourceEnvName,
			}

			// Check if the environment matches the filters
			ok, err := utils.MatchesFilters(formattedEnvironment, filterMaps)
			if err != nil {
				utils.PrintError(err)
				return err
			}
			if ok {
				environments = append(environments, formattedEnvironment)
			}
		}
	}

	var jsonData []string
	for _, environment := range environments {
		data, err := json.MarshalIndent(environment, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		jsonData = append(jsonData, string(data))
	}

	if len(jsonData) == 0 {
		utils.PrintInfo("No environments found.")
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
