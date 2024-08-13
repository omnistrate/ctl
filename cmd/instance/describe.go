package instance

import (
	"encoding/json"
	"fmt"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe the instance deployment
omnistrate instance describe --instance-id=instance-12345`
)

var describeCmd = &cobra.Command{
	Use:          "describe --instance-id=[instance-id]",
	Short:        "Describe a instance deployment your service.",
	Long:         `This command helps you describe the instance for your service.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Flags().String("instance-id", "", "Instance ID")

	err := describeCmd.MarkFlagRequired("instance-id")
	if err != nil {
		return
	}
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Get flags

	instanceId, _ := cmd.Flags().GetString("instance-id")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Check if the instance exists
	searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", instanceId))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var found bool
	var serviceId, environmentId string
	for _, instance := range searchRes.ResourceInstanceResults {
		if instance.ID == instanceId {
			serviceId = string(instance.ServiceID)
			environmentId = string(instance.ServiceEnvironmentID)
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceId)
		utils.PrintError(err)
		return nil
	}

	var instance *inventoryapi.ResourceInstance
	instance, err = dataaccess.DescribeInstance(token, serviceId, environmentId, instanceId)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	data, err := json.MarshalIndent(instance, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
