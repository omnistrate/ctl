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
omnistrate instance describe --service-id=s-12345 --service-environment-id=se-12345 --instance-id=instance-12345`
)

var describeCmd = &cobra.Command{
	Use:          "describe --service-id=[service-id] --service-environment-id=[service-environment-id] --instance-id=[instance-id]",
	Short:        "Describe a instance deployment your service.",
	Long:         `This command helps you describe the instance for your service.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Flags().String("service-id", "", "Service ID")
	describeCmd.Flags().String("service-environment-id", "", "Service Environment ID")
	describeCmd.Flags().String("instance-id", "", "Instance ID")

	err := describeCmd.MarkFlagRequired("service-id")
	if err != nil {
		return
	}

	err = describeCmd.MarkFlagRequired("service-environment-id")
	if err != nil {
		return
	}

	err = describeCmd.MarkFlagRequired("instance-id")
	if err != nil {
		return
	}
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Get flags
	serviceId, _ := cmd.Flags().GetString("service-id")
	serviceEnvironmentId, _ := cmd.Flags().GetString("service-environment-id")
	instanceId, _ := cmd.Flags().GetString("instance-id")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var instance *inventoryapi.ResourceInstance
	instance, err = dataaccess.DescribeInstance(token, serviceId, serviceEnvironmentId, instanceId)
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
