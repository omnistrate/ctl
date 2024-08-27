package serviceplan

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	tierversionsetapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/tier_version_set_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	describeVersionExample = `  # Describe a service plan version
  omctl service-plan describe-version [service-name] [plan-name] --version [version]

  # Describe a service plan version by ID instead of name
  omctl service-plan describe-version --service-id [service-id] --plan-id [plan-id] --version [version]`

	defaultDescribeVersionOutput = "json"
)

var describeVersionCmd = &cobra.Command{
	Use:          "describe-version [service-name] [plan-name] [flags]",
	Short:        "Describe a service plan version",
	Long:         `This command helps you describe a service plan version in your service.`,
	Example:      describeVersionExample,
	RunE:         runDescribeVersion,
	SilenceUsage: true,
}

func init() {
	describeVersionCmd.Flags().StringP("version", "v", "", "Service plan version (latest|preferred|1.0 etc.)")
	describeVersionCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	describeVersionCmd.Flags().StringP("plan-id", "", "", "Environment ID. Required if plan name is not provided")

	err := describeVersionCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

func runDescribeVersion(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	serviceID, _ := cmd.Flags().GetString("service-id")
	planID, _ := cmd.Flags().GetString("plan-id")
	version, _ := cmd.Flags().GetString("version")
	output := defaultDescribeVersionOutput

	// Validate input arguments
	if err := validateDescribeVersionArguments(args, serviceID, planID); err != nil {
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
	serviceID, serviceName, planID, _, environment, err := getServicePlan(token, serviceID, serviceName, planID, planName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get the target version
	version, err = getTargetVersion(token, serviceID, planID, version)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe the version set
	servicePlan, err := dataaccess.DescribeVersionSet(token, serviceID, planID, version)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Format the service plan details
	formattedServicePlan, err := formatServicePlanVersionDetails(token, serviceName, planName, environment, servicePlan)
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

func validateDescribeVersionArguments(args []string, serviceID, planID string) error {
	if len(args) == 0 && (serviceID == "" || planID == "") {
		return fmt.Errorf("please provide the service name and service plan name or the service ID and service plan ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [plan-name]", strings.Join(args, " "))
	}
	return nil
}

func formatServicePlanVersionDetails(token, serviceName, planName, environment string, versionSet *tierversionsetapi.TierVersionSet) (model.ServicePlanVersionDetails, error) {
	// Get resource details
	var resources []model.Resource
	for _, versionSetResource := range versionSet.Resources {
		// Get resource details
		desRes, err := dataaccess.DescribeResource(token, string(versionSet.ServiceID), string(versionSetResource.ID), commonutils.ToPtr(string(versionSet.ProductTierID)), &versionSet.Version)
		if err != nil {
			return model.ServicePlanVersionDetails{}, err
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

	formattedServicePlan := model.ServicePlanVersionDetails{
		PlanID:             string(versionSet.ProductTierID),
		PlanName:           planName,
		ServiceID:          string(versionSet.ServiceID),
		ServiceName:        serviceName,
		Environment:        environment,
		Version:            versionSet.Version,
		ReleaseDescription: utils.GetStrValue(versionSet.Name),
		VersionSetStatus:   versionSet.Status,
		EnabledFeatures:    versionSet.EnabledFeatures,
		Resources:          resources,
	}

	return formattedServicePlan, nil
}
