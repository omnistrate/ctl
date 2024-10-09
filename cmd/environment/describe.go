package environment

import (
	"fmt"
	"strings"

	"github.com/chelnak/ysmrr"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe environment
omctl environment describe [service-name] [environment-name]

# Describe environment by ID instead of name
omctl environment describe --service-id=[service-id] --environment-id=[environment-id]`
)

var describeCmd = &cobra.Command{
	Use:          "describe [service-name] [environment-name] [flags]",
	Short:        "Describe a Service Environment",
	Long:         `This command helps you get details of a service environment from your service. You can find details like SaaS portal status, SaaS portal URL, and promote status, etc.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	describeCmd.Flags().StringP("environment-id", "", "", "Environment ID. Required if environment name is not provided")
	describeCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported.") // Override inherited flag
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	serviceID, _ := cmd.Flags().GetString("service-id")
	environmentID, _ := cmd.Flags().GetString("environment-id")
	output, _ := cmd.Flags().GetString("output")

	// Validate input arguments
	if err := validateDescribeArguments(args, serviceID, environmentID, output); err != nil {
		utils.PrintError(err)
		return err
	}

	// Set service and environment names if provided in args
	var serviceName, environmentName string
	if len(args) == 2 {
		serviceName, environmentName = args[0], args[1]
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
		spinner = sm.AddSpinner("Describing environment...")
		sm.Start()
	}

	// Retrieve service and environment details
	serviceID, serviceName, environmentID, _, err = getServiceEnvironment(token, serviceID, serviceName, environmentID, environmentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe the environment
	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceID, environmentID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get the source environment name
	sourceEnvName := ""
	if environment.SourceEnvironmentID != nil {
		sourceEnv, err := dataaccess.DescribeServiceEnvironment(token, serviceID, string(*environment.SourceEnvironmentID))
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		sourceEnvName = sourceEnv.Name
	}

	// Format the environment details
	formattedEnvironment := formatEnvironmentDetails(token, serviceID, serviceName, sourceEnvName, environment)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Handle output based on format
	if spinner != nil {
		spinner.UpdateMessage("Successfully retrieved environment details")
		spinner.Complete()
		sm.Stop()
	}

	if err = utils.PrintTextTableJsonOutput(output, formattedEnvironment); err != nil {
		return err
	}

	return nil
}

// Helper functions

func validateDescribeArguments(args []string, serviceID, environmentID, output string) error {
	if len(args) == 0 && (serviceID == "" || environmentID == "") {
		return fmt.Errorf("please provide the service name and environment name or the service ID and environment ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [environment-name]", strings.Join(args, " "))
	}
	if output != "json" {
		return fmt.Errorf("only json output is supported")
	}
	return nil
}

func formatEnvironmentDetails(token, serviceID, serviceName, sourceEnvName string, environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) model.DetailedEnvironment {
	formattedEnvironment := model.DetailedEnvironment{
		EnvironmentID:    string(environment.ID),
		EnvironmentName:  environment.Name,
		EnvironmentType:  string(environment.Type),
		ServiceID:        string(environment.ServiceID),
		SourceEnvName:    sourceEnvName,
		ServiceName:      serviceName,
		SaaSPortalStatus: getSaaSPortalStatus(environment),
		SaaSPortalURL:    getSaaSPortalURL(environment),
		PromoteStatus:    getPromoteStatus(token, serviceID, environment),
	}

	return formattedEnvironment
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
