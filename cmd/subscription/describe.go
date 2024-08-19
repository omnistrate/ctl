package subscription

import (
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
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
	for _, subscription := range searchRes.SubscriptionResults {
		if subscription.ID == subscriptionId {
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("%s not found. Please check the subscription ID and try again", subscriptionId)
		utils.PrintError(err)
		return nil
	}

	subscription := searchRes.SubscriptionResults[0]
	envType := ""
	if subscription.ServiceEnvironmentType != nil {
		envType = string(*subscription.ServiceEnvironmentType)
	}

	formattedSubscription := model.Subscription{
		SubscriptionID:         subscription.ID,
		ServiceID:              string(subscription.ServiceID),
		ServiceName:            subscription.ServiceName,
		PlanID:                 string(subscription.ProductTierID),
		PlanName:               subscription.ServicePlanName,
		Environment:            envType,
		SubscriptionOwnerName:  subscription.RootUserName,
		SubscriptionOwnerEmail: subscription.RootUserEmail,
		Status:                 string(subscription.Status),
	}

	data, err := json.MarshalIndent(formattedSubscription, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
