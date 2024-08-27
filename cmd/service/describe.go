package service

import (
	"fmt"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	describeExample = `  # Describe service with name
  omctl service describe [service-name]

  # Describe service with ID
  omctl service describe --id=[service-ID]`
)

var describeCmd = &cobra.Command{
	Use:          "describe [flags]",
	Short:        "Describe a service",
	Long:         "This command helps you describe a service using its name or ID.",
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	describeCmd.Flags().String("id", "", "Service ID")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	var name string
	if len(args) > 0 {
		name = args[0]
	}

	// Retrieve flags
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate input args
	err = validateDescribeInput(name, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Check if service exists
	id, name, err = getService(token, name, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Describe service
	var service *serviceapi.DescribeServiceResult
	service, err = dataaccess.DescribeService(token, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	err = utils.PrintTextTableJsonOutput("json", service)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}

// Helper functions

func getService(token, serviceNameArg, serviceIDArg string) (serviceID, serviceName string, err error) {
	if serviceNameArg != "" {
		var searchRes *inventoryapi.SearchInventoryResult
		searchRes, err = dataaccess.SearchInventory(token, fmt.Sprintf("service:%s", serviceNameArg))
		if err != nil {
			return
		}

		for _, service := range searchRes.ServiceResults {
			if service.Name == serviceNameArg {
				serviceID = service.ID
				serviceName = service.Name
				return
			}
		}
	} else {
		var searchRes *inventoryapi.SearchInventoryResult
		searchRes, err = dataaccess.SearchInventory(token, fmt.Sprintf("service:%s", serviceIDArg))
		if err != nil {
			return
		}

		for _, service := range searchRes.ServiceResults {
			if service.ID == serviceIDArg {
				serviceID = service.ID
				serviceName = service.Name
				return
			}
		}
	}

	return
}

func validateDescribeInput(serviceNameArg, serviceIDArg string) error {
	if serviceNameArg == "" && serviceIDArg == "" {
		return errors.New("service name or ID must be provided")
	}

	if serviceNameArg != "" && serviceIDArg != "" {
		return errors.New("only one of service name or ID can be provided")
	}

	return nil
}
