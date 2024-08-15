package instance

import (
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/table"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

const (
	listExample = `# List instances of the service postgres in the prod environment
omnistrate instance list --output=table --filters "service:postgres,environment:prod"`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all instances in your account",
	Long: `This command helps you list all the instances in your account.
You can filter for specific instances by using the filter flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringP("filters", "f", "", "Filter instances by a specific criteria")
	listCmd.Flags().Bool("truncate-names", false, "Truncate long names in the output")
}

type Instance struct {
	InstanceID    string `json:"instance_id"`
	Service       string `json:"service"`
	Environment   string `json:"environment"`
	PlanName      string `json:"plan_name"`
	PlanVersion   string `json:"plan_version"`
	Resource      string `json:"resource"`
	CloudProvider string `json:"cloud_provider"`
	Region        string `json:"region"`
	Status        string `json:"status"`
}

func runList(cmd *cobra.Command, args []string) error {
	// Get flags
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	truncateNames, err := cmd.Flags().GetBool("truncate-names")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	// TODO: Implement filters

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Check if the instance exists
	searchRes, err := dataaccess.SearchInventory(token, "resourceinstance:%s")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	instances := make([]Instance, 0)
	for _, instance := range searchRes.ResourceInstanceResults {
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

		instances = append(instances, Instance{
			InstanceID:    instance.ID,
			Service:       serviceName,
			Environment:   envType,
			PlanName:      planName,
			PlanVersion:   planVersion,
			Resource:      instance.ResourceName,
			CloudProvider: string(instance.CloudProvider),
			Region:        instance.RegionCode,
			Status:        string(instance.Status),
		})
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
		printTable(instances)
	case "table":
		var tableWriter *table.Table
		if tableWriter, err = table.NewTableFromJSONTemplate(json.RawMessage(jsonData[0])); err != nil {
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

func printTable(res []Instance) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	_, err := fmt.Fprintln(w, "Instance ID\tService\tEnvironment\tPlan Name\tPlan Version\tResource\tCloud Provider\tRegion\tStatus")
	if err != nil {
		return
	}

	for _, r := range res {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			r.InstanceID,
			r.Service,
			r.Environment,
			r.PlanName,
			r.PlanVersion,
			r.Resource,
			r.CloudProvider,
			r.Region,
			r.Status,
		)
		if err != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
