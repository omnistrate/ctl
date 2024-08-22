package serviceplan

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
	describeExample = `# Describe service plan
omnistrate service-plan describe [service-name] [plan-name]

# Describe service plan by ID instead of name
omnistrate service-plan describe --service-id [service-id] --plan-id [plan-id]`

	defaultDescribeOutput = "json"
)

var describeCmd = &cobra.Command{
	Use:          "describe [service-name] [plan-name] [flags]",
	Short:        "Describe a service plan",
	Long:         `This command helps you describe a service plan in your service.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Flags().StringP("version", "v", "", "Service plan version (latest|preferred|1.0 etc.)")
	describeCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	describeCmd.Flags().StringP("plan-id", "", "", "Environment ID. Required if plan name is not provided")

	err := describeCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	serviceId, _ := cmd.Flags().GetString("service-id")
	planId, _ := cmd.Flags().GetString("plan-id")
	version, _ := cmd.Flags().GetString("version")
	output := defaultDescribeOutput

	// Validate input arguments
	if err := validateDescribeArguments(args, serviceId, planId); err != nil {
		utils.PrintError(err)
		return err
	}

	// Set service and service plan names if provided in args
	var serviceName, planName string
	if len(args) == 2 {
		serviceName, planName = args[0], args[1]
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
		spinner = sm.AddSpinner("Describing service plan...")
		sm.Start()
	}

	// Retrieve service and service plan details
	services, err := dataaccess.ListServices(token)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	serviceId, serviceName, planId, _, err = getServicePlan(services, serviceId, serviceName, planId, planName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get the target version
	version, err = getTargetVersion(token, serviceId, planId, version)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe the version set
	servicePlan, err := dataaccess.DescribeVersionSet(token, serviceId, planId, version)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Format the service plan details
	formattedServicePlan, err := formatServicePlanDetails(token, serviceId, serviceName, servicePlan)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Handle output based on format
	if spinner != nil {
		spinner.UpdateMessage("Service plan description retrieved successfully")
		spinner.Complete()
		sm.Stop()
	}

	if err = utils.PrintTextTableJsonOutput(output, formattedServicePlan); err != nil {
		return err
	}

	return nil
}

// Helper functions

func validateDescribeArguments(args []string, serviceId, planId string) error {
	if len(args) == 0 && (serviceId == "" || planId == "") {
		return fmt.Errorf("please provide the service name and service plan name or the service ID and service plan ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [plan-name]", strings.Join(args, " "))
	}
	return nil
}

func formatServicePlanDetails(token, serviceId, serviceName string, environment *serviceenvironmentapi.DescribeServicePlanResult) (string, error) {
	// Example of formatting environment details
	formattedServicePlan := model.ServicePlanDetails{
		EnvironmentID:    string(environment.ID),
		EnvironmentName:  environment.Name,
		EnvironmentType:  string(environment.Type),
		ServiceID:        string(environment.ServiceID),
		ServiceName:      serviceName,
		SaaSPortalStatus: getSaaSPortalStatus(environment),
		SaaSPortalURL:    getSaaSPortalURL(environment),
		PromoteStatus:    getPromoteStatus(token, serviceId, environment),
	}

	data, err := json.MarshalIndent(formattedServicePlan, "", "    ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}
