package serviceplan

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	producttierapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/product_tier_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
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
	describeCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	describeCmd.Flags().StringP("plan-id", "", "", "Environment ID. Required if plan name is not provided")
}

func runDescribe(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	serviceId, _ := cmd.Flags().GetString("service-id")
	planId, _ := cmd.Flags().GetString("plan-id")
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

	// Check if the service plan exists
	serviceId, serviceName, planId, _, environment, err := getServicePlan(token, serviceId, serviceName, planId, planName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe the service plan
	servicePlan, err := dataaccess.DescribeProductTier(token, serviceId, planId)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Format the service plan details
	formattedServicePlan, err := formatServicePlanDetails(token, serviceName, planName, environment, servicePlan)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Handle output based on format
	utils.HandleSpinnerSuccess(spinner, sm, "Service plan details retrieved successfully")

	// Marshal data
	data, err := json.MarshalIndent(formattedServicePlan, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if err = utils.PrintTextTableJsonOutput(output, string(data)); err != nil {
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

func formatServicePlanDetails(token, serviceName, planName, environment string, productTier *producttierapi.DescribeProductTierResult) (model.ServicePlanDetails, error) {
	// Get service model
	serviceModel, err := dataaccess.DescribeServiceModel(token, string(productTier.ServiceID), string(productTier.ServiceModelID))
	if err != nil {
		return model.ServicePlanDetails{}, err
	}

	// Get resource details
	var resources []model.Resource
	for resourceID := range productTier.APIGroups {
		// Get resource details
		desRes, err := dataaccess.DescribeResource(token, string(productTier.ServiceID), string(resourceID), commonutils.ToPtr(string(productTier.ID)), nil)
		if err != nil {
			return model.ServicePlanDetails{}, err
		}
		resource := model.Resource{
			ResourceID:   string(desRes.ID),
			ResourceName: desRes.Name,
			ResourceType: string(desRes.ResourceType),
		}

		if desRes.ActionHooks != nil {
			resource.ActionHooks = desRes.ActionHooks
		}
		if desRes.AdditionalSecurityContext != nil {
			resource.AdditionalSecurityContext = desRes.AdditionalSecurityContext
		}
		if desRes.BackupConfiguration != nil {
			resource.BackupConfiguration = desRes.BackupConfiguration
		}
		if desRes.Capabilities != nil {
			resource.Capabilities = desRes.Capabilities
		}
		if desRes.CustomLabels != nil {
			resource.CustomLabels = desRes.CustomLabels
		}
		if desRes.CustomSysCTLs != nil {
			resource.CustomSysCTLs = desRes.CustomSysCTLs
		}
		if desRes.CustomULimits != nil {
			resource.CustomULimits = desRes.CustomULimits
		}
		if desRes.Dependencies != nil {
			resource.Dependencies = desRes.Dependencies
		}
		if desRes.EnvironmentVariables != nil {
			resource.EnvironmentVariables = desRes.EnvironmentVariables
		}
		if desRes.FileSystemConfiguration != nil {
			resource.FileSystemConfiguration = desRes.FileSystemConfiguration
		}
		if desRes.HelmChartConfiguration != nil {
			resource.HelmChartConfiguration = desRes.HelmChartConfiguration
		}
		if desRes.KustomizeConfiguration != nil {
			resource.KustomizeConfiguration = desRes.KustomizeConfiguration
		}
		if desRes.L4LoadBalancerConfiguration != nil {
			resource.L4LoadBalancerConfiguration = desRes.L4LoadBalancerConfiguration
		}
		if desRes.L7LoadBalancerConfiguration != nil {
			resource.L7LoadBalancerConfiguration = desRes.L7LoadBalancerConfiguration
		}
		if desRes.OperatorCRDConfiguration != nil {
			resource.OperatorCRDConfiguration = desRes.OperatorCRDConfiguration
		}

		resources = append(resources, resource)
	}

	formattedServicePlan := model.ServicePlanDetails{
		PlanID:          string(productTier.ID),
		PlanName:        planName,
		ServiceID:       string(productTier.ServiceID),
		ServiceName:     serviceName,
		Environment:     environment,
		DeploymentType:  string(productTier.TierType),
		TenancyType:     string(serviceModel.ModelType),
		EnabledFeatures: productTier.EnabledFeatures,
		Resources:       resources,
	}

	return formattedServicePlan, nil
}
