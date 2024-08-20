package environment

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	deleteExample = `# Delete environment
omnistrate environment delete [service-name] [environment-name]

# Delete environment by ID instead of name
omnistrate environment delete --service-id [service-id] --environment-id [environment-id]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [service-name] [environment-name] [flags]",
	Short:        "Delete a environment",
	Long:         `This command helps you delete a environment in your service.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	deleteCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	deleteCmd.Flags().StringP("environment-id", "", "", "Environment ID. Required if environment name is not provided")
}

func runDelete(cmd *cobra.Command, args []string) error {
	// Get flags
	output, _ := cmd.Flags().GetString("output")
	serviceId, _ := cmd.Flags().GetString("service-id")
	environmentId, _ := cmd.Flags().GetString("environment-id")

	if len(args) == 0 && (serviceId == "" || environmentId == "") {
		err := fmt.Errorf("please provide the service name and environment name or the service ID and environment ID")
		utils.PrintError(err)
		return err
	}

	if len(args) > 0 && len(args) != 2 {
		err := fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [environment-name]", strings.Join(args, " "))
		utils.PrintError(err)
		return err
	}

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Deleting environment..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if the environment exists
	services, err := dataaccess.ListServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	serviceEnvironmentsMap := make(map[string]map[string]bool)
	for _, service := range services.Services {
		if string(service.ID) == serviceId || (len(args) == 2 && strings.EqualFold(service.Name, args[0])) {
			for _, environment := range service.ServiceEnvironments {
				if string(environment.ID) == environmentId || (len(args) == 2 && strings.EqualFold(environment.Name, args[1])) {
					if _, ok := serviceEnvironmentsMap[string(service.ID)]; !ok {
						serviceEnvironmentsMap[string(service.ID)] = make(map[string]bool)
					}
					serviceEnvironmentsMap[string(service.ID)][string(environment.ID)] = true
					serviceId = string(service.ID)
					environmentId = string(environment.ID)
				}
			}
		}
	}
	if len(serviceEnvironmentsMap) == 0 {
		err = errors.New("environment not found. Please check the input values and try again")
		utils.PrintError(err)
		return err
	}
	if len(serviceEnvironmentsMap) > 1 || len(serviceEnvironmentsMap[serviceId]) > 1 {
		err = errors.New("multiple environments found. Please provide the service ID and environment ID instead of the names")
		utils.PrintError(err)
		return err
	}

	err = dataaccess.DeleteServiceEnvironment(token, serviceId, environmentId)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if output != "json" {
		spinner.UpdateMessage("Successfully deleted environment")
		spinner.Complete()
		sm.Stop()
	}

	return nil
}
