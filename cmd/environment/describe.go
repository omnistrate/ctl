package environment

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	describeExample = `# Describe environment
omnistrate environment describe [service-name] [environment-name]

# Describe environment by ID instead of name
omnistrate environment describe --service-id [service-id] --environment-id [environment-id]`

	defaultDescribeOutput = "json"
)

var describeCmd = &cobra.Command{
	Use:          "describe [service-name] [environment-name] [flags]",
	Short:        "Describe a environment",
	Long:         `This command helps you describe a environment in your service.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	describeCmd.Flags().StringP("environment-id", "", "", "Environment ID. Required if environment name is not provided")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer cleanUpDescribeFlagsAndArgs(cmd, &args)

	// Retrieve flags
	serviceId, _ := cmd.Flags().GetString("service-id")
	environmentId, _ := cmd.Flags().GetString("environment-id")
	output := defaultDescribeOutput

	// Validate input arguments
	if err := validateDescribeArguments(args, serviceId, environmentId); err != nil {
		utils.PrintError(err)
		return err
	}

	// Set service and environment names if provided in args
	var serviceName, environmentName string
	if len(args) == 2 {
		serviceName, environmentName = args[0], args[1]
	}

	// Validate user login
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		spinner = sm.AddSpinner("Describing environment...")
		sm.Start()
	}

	// Retrieve service and environment details
	services, err := dataaccess.ListServices(token)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	serviceId, serviceName, environmentId, environmentName, err = getServiceEnvironment(services, serviceId, serviceName, environmentId, environmentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe the environment
	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceId, environmentId)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Format the environment details
	formattedEnvironment, err := formatEnvironmentDetails(token, serviceId, serviceName, environment)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Handle output based on format
	if spinner != nil {
		spinner.UpdateMessage("Environment description retrieved successfully")
		spinner.Complete()
		sm.Stop()
	}

	if err = utils.PrintTextTableJsonOutput(output, formattedEnvironment); err != nil {
		return err
	}

	return nil
}

// Helper functions

func validateDescribeArguments(args []string, serviceId, environmentId string) error {
	if len(args) == 0 && (serviceId == "" || environmentId == "") {
		return fmt.Errorf("please provide the service name and environment name or the service ID and environment ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [environment-name]", strings.Join(args, " "))
	}
	return nil
}

func formatEnvironmentDetails(token, serviceId, serviceName string, environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) (string, error) {
	// Example of formatting environment details
	formattedEnvironment := model.DetailedEnvironment{
		EnvironmentID:    string(environment.ID),
		EnvironmentName:  environment.Name,
		EnvironmentType:  string(environment.Type),
		ServiceID:        string(environment.ServiceID),
		ServiceName:      serviceName,
		SaaSPortalStatus: getSaaSPortalStatus(environment),
		SaaSPortalURL:    getSaaSPortalURL(environment),
		PromoteStatus:    getPromoteStatus(token, serviceId, environment),
	}

	data, err := json.MarshalIndent(formattedEnvironment, "", "    ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getSaaSPortalStatus(environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) string {
	if environment.SaasPortalStatus != nil {
		return string(*environment.SaasPortalStatus)
	}
	return ""
}

func getSaaSPortalURL(environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) string {
	if environment.SaasPortalURL != nil {
		return *environment.SaasPortalURL
	}
	return ""
}

func cleanUpDescribeFlagsAndArgs(cmd *cobra.Command, args *[]string) {
	// Clean up flags
	_ = cmd.Flags().Set("service-id", "")
	_ = cmd.Flags().Set("environment-id", "")

	// Clean up arguments by resetting the slice to nil or an empty slice
	*args = nil
}
