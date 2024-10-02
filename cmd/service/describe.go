package service

import (
	"fmt"
	"strings"

	"github.com/chelnak/ysmrr"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe service with name
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
	describeCmd.Args = cobra.MaximumNArgs(1) // Require at most one argument

	describeCmd.Flags().String("id", "", "Service ID")
	describeCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported.") // Override inherited flag
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

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
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate input args
	err = validateDescribeArguments(name, id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Describing service..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if service exists
	id, _, err = getService(token, name, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe service
	var service *serviceapi.DescribeServiceResult
	service, err = dataaccess.DescribeService(token, id)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully described service")

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
	count := 0
	if serviceNameArg != "" {
		var searchRes *inventoryapi.SearchInventoryResult
		searchRes, err = dataaccess.SearchInventory(token, fmt.Sprintf("service:%s", serviceNameArg))
		if err != nil {
			return
		}

		for _, service := range searchRes.ServiceResults {
			if strings.EqualFold(service.Name, serviceNameArg) {
				serviceID = service.ID
				serviceName = service.Name
				count++
			}
		}
	} else {
		var searchRes *inventoryapi.SearchInventoryResult
		searchRes, err = dataaccess.SearchInventory(token, fmt.Sprintf("service:%s", serviceIDArg))
		if err != nil {
			return
		}

		for _, service := range searchRes.ServiceResults {
			if strings.EqualFold(service.ID, serviceIDArg) {
				serviceID = service.ID
				serviceName = service.Name
				count++
			}
		}
	}

	if count == 0 {
		err = errors.New("service not found")
		return
	}

	if count > 1 {
		err = errors.New("multiple services found with the same name. Please provide the service ID")
	}

	return
}

func validateDescribeArguments(serviceNameArg, serviceIDArg string) error {
	if serviceNameArg == "" && serviceIDArg == "" {
		return errors.New("service name or ID must be provided")
	}

	if serviceNameArg != "" && serviceIDArg != "" {
		return errors.New("only one of service name or ID can be provided")
	}

	return nil
}
