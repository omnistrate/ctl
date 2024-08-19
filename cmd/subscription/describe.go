package subscription

import (
	"encoding/json"
	"fmt"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe subscription
omnistrate subscription describe subscription-abcd1234`
)

var describeCmd = &cobra.Command{
	Use:          "describe [subscription-id]",
	Short:        "Describe a customer subscription to your service",
	Long:         `This command helps you describe a customer subscription to your service.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Get flags
	subscriptionId := args[0]

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Check if the subscription exists
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("subscription:%s", subscriptionId))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var found bool
	var serviceId, environmentId string
	for _, subscription := range searchRes.SubscriptionResults {
		if subscription.ID == subscriptionId {
			serviceId = string(subscription.ServiceID)
			environmentId = string(subscription.ServiceEnvironmentID)
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("%s not found. Please check the subscription ID and try again", subscriptionId)
		utils.PrintError(err)
		return nil
	}

	var subscription *inventoryapi.FleetDescribeSubscriptionResult
	subscription, err = dataaccess.DescribeSubscription(token, serviceId, environmentId, subscriptionId)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	data, err := json.MarshalIndent(subscription, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
