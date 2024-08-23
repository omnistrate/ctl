package subscription

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
	listExample = `  # List subscriptions of the service postgres and mysql in the prod environment
  omctl subscription list -o=table -f="service_name:postgres,environment:PROD" -f="service:mysql,environment:PROD"`
	defaultMaxNameLength = 30 // Maximum length of the name column in the table
)

var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List customer subscriptions to your services",
	Long: `This command helps you list customer subscriptions to your services.
You can filter for specific subscriptions by using the filter flag.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	listCmd.Flags().StringArrayP("filter", "f", []string{}, "Filter to apply to the list of subscriptions. E.g.: key1:value1,key2:value2, which filters subscriptions where key1 equals value1 and key2 equals value2. Allow use of multiple filters to form the logical OR operation. Supported keys: "+strings.Join(utils.GetSupportedFilterKeys(model.Subscription{}), ",")+". Check the examples for more details.")
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
	filterMaps, err := utils.ParseFilters(filters, utils.GetSupportedFilterKeys(model.Subscription{}))
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

	// Get all subscriptions
	searchRes, err := dataaccess.SearchInventory(token, "subscription:s")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	subscriptions := make([]model.Subscription, 0)
	for _, subscription := range searchRes.SubscriptionResults {
		if subscription == nil {
			continue
		}
		serviceName := subscription.ServiceName
		planName := subscription.ServicePlanName
		if truncateNames {
			serviceName = utils.TruncateString(serviceName, defaultMaxNameLength)
			planName = utils.TruncateString(planName, defaultMaxNameLength)
		}
		formattedSubscription := model.Subscription{
			SubscriptionID:         subscription.ID,
			ServiceID:              string(subscription.ServiceID),
			ServiceName:            serviceName,
			PlanID:                 string(subscription.ProductTierID),
			PlanName:               planName,
			Environment:            subscription.ServiceEnvironmentName,
			SubscriptionOwnerName:  subscription.RootUserName,
			SubscriptionOwnerEmail: subscription.RootUserEmail,
			Status:                 string(subscription.Status),
		}

		// Check if the subscription matches the filters
		ok, err := utils.MatchesFilters(formattedSubscription, filterMaps)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		if ok {
			subscriptions = append(subscriptions, formattedSubscription)
		}
	}

	var jsonData []string
	for _, subscription := range subscriptions {
		data, err := json.MarshalIndent(subscription, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		jsonData = append(jsonData, string(data))
	}

	if len(jsonData) == 0 {
		utils.PrintInfo("No subscriptions found.")
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
