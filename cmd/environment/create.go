package environment

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"slices"
	"strings"
)

const (
	createExample = `# Create environment
omnistrate environment create [service-name] [environment-name] --type [type] --source [source]

# Create environment by ID instead of name
omnistrate environment create [environment-name] --service-id [service-id] --type [type] --source [source]`
)

var EnvironmentID string

var createCmd = &cobra.Command{
	Use:          "create [service-name] [environment-name] [flags]",
	Short:        "Create a environment",
	Long:         `This command helps you create a environment in your service.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Flags().StringP("description", "", "", "Environment description")
	createCmd.Flags().StringP("type", "", "", "Type of environment. Valid options include: 'dev', 'prod', 'qa', 'canary', 'staging', 'private'")
	createCmd.Flags().StringP("source", "", "", "Source environment name")
	createCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	createCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")

	err := createCmd.MarkFlagRequired("type")
	if err != nil {
		return
	}

	createCmd.Args = cobra.MinimumNArgs(1)
	createCmd.Args = cobra.MaximumNArgs(2)
}

func runCreate(cmd *cobra.Command, args []string) error {
	defer cleanUpCreateFlagsAndArgs(cmd, &args)

	// Get flags
	description, _ := cmd.Flags().GetString("description")
	envType, _ := cmd.Flags().GetString("type")
	sourceEnvName, _ := cmd.Flags().GetString("source")
	output, _ := cmd.Flags().GetString("output")
	serviceId, _ := cmd.Flags().GetString("service-id")
	envName := args[len(args)-1]

	// Set service name if provided in args
	var serviceName string
	if len(args) == 2 {
		serviceName = args[0]
	}

	// Validate input arguments
	if err := validateCreateArguments(args, serviceId); err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate environment type
	if err := validateEnvironmentType(envType); err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user is logged in
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
		spinner = sm.AddSpinner("Creating environment...")
		sm.Start()
	}

	// Check if the service exists
	serviceId, serviceName, err = getService(token, serviceId, serviceName)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Check if source environment exists
	sourceEnvID, err := getSourceEnvironmentID(token, serviceId, sourceEnvName)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Create the environment
	if description == "" {
		description = fmt.Sprintf("%s environment for service %s", envType, serviceName)
	}

	visibility := getVisibility(envType)
	defaultDeploymentConfigID, err := dataaccess.GetDefaultDeploymentConfigID(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	publicKeyPtr := getPublicKeyPtr(visibility)
	environmentID, err := createEnvironment(token, envName, description, serviceId, envType, visibility, string(defaultDeploymentConfigID), sourceEnvID, publicKeyPtr)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully created environment")

	// Describe the environment
	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceId, string(environmentID))
	if err != nil {
		return err
	}

	// Format and print the environment
	formattedEnvironment, err := formatDetailedEnvironment(token, serviceId, serviceName, sourceEnvName, environment)
	if err != nil {
		return err
	}
	err = utils.PrintTextTableJsonOutput(output, formattedEnvironment)
	if err != nil {
		return err
	}

	EnvironmentID = string(environmentID)

	return nil
}

// Helper functions

func validateCreateArguments(args []string, serviceId string) error {
	if len(args) == 1 && serviceId == "" {
		return fmt.Errorf("please provide service name or id")
	}
	if len(args) == 2 && serviceId != "" {
		return fmt.Errorf("please provide either service name or id, not both")
	}
	return nil
}

func validateEnvironmentType(envType string) error {
	if !slices.Contains([]string{"dev", "qa", "staging", "canary", "prod", "private"}, strings.ToLower(envType)) {
		return fmt.Errorf("invalid environment type: %s", envType)
	}
	return nil
}

func getService(token, serviceIdArg, serviceNameArg string) (serviceId string, serviceName string, err error) {
	searchRes, err := dataaccess.SearchInventory(token, "service:s")
	if err != nil {
		return "", "", err
	}

	for _, service := range searchRes.ServiceResults {
		if serviceIdArg == service.ID || strings.EqualFold(service.Name, serviceNameArg) {
			return service.ID, service.Name, nil
		}
	}

	return "", "", errors.New("service not found")
}

func getSourceEnvironmentID(token, serviceId, sourceEnvName string) (string, error) {
	if sourceEnvName == "" {
		return "", nil
	}

	describeServiceRes, err := dataaccess.DescribeService(token, serviceId)
	if err != nil {
		return "", err
	}

	for _, env := range describeServiceRes.ServiceEnvironments {
		if strings.EqualFold(env.Name, sourceEnvName) {
			return string(env.ID), nil
		}
	}

	return "", errors.New("source environment not found. Please provide a valid source environment name")
}

func getVisibility(envType string) string {
	switch strings.ToLower(envType) {
	case "dev", "qa", "staging", "canary", "private":
		return "PRIVATE"
	case "prod":
		return "PUBLIC"
	default:
		return ""
	}
}

func getPublicKeyPtr(visibility string) *string {
	if visibility == "PRIVATE" {
		return commonutils.ToPtr(utils.GetDefaultServiceAuthPublicKey())
	}
	return nil
}

func createEnvironment(token, envName, description, serviceId, envType, visibility, defaultDeploymentConfigID, sourceEnvID string, publicKeyPtr *string) (serviceenvironmentapi.ServiceEnvironmentID, error) {
	request := serviceenvironmentapi.CreateServiceEnvironmentRequest{
		Name:                    envName,
		Description:             description,
		ServiceID:               serviceenvironmentapi.ServiceID(serviceId),
		Visibility:              serviceenvironmentapi.ServiceVisibility(visibility),
		Type:                    commonutils.ToPtr(serviceenvironmentapi.EnvironmentType(envType)),
		ServiceAuthPublicKey:    publicKeyPtr,
		DeploymentConfigID:      serviceenvironmentapi.DeploymentConfigID(defaultDeploymentConfigID),
		AutoApproveSubscription: commonutils.ToPtr(true),
		SourceEnvironmentID:     (*serviceenvironmentapi.ServiceEnvironmentID)(commonutils.ToPtr(sourceEnvID)),
	}

	return dataaccess.CreateServiceEnvironment(token, request)
}

func formatDetailedEnvironment(token, serviceId, serviceName, sourceEnvName string, environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) (string, error) {
	saasPortalStatus := ""
	if environment.SaasPortalStatus != nil {
		saasPortalStatus = string(*environment.SaasPortalStatus)
	}
	saasPortalURL := ""
	if environment.SaasPortalURL != nil {
		saasPortalURL = *environment.SaasPortalURL
	}

	promoteStatus := getPromoteStatus(token, serviceId, environment)

	formattedEnvironment := model.DetailedEnvironment{
		EnvironmentID:    string(environment.ID),
		EnvironmentName:  environment.Name,
		EnvironmentType:  string(environment.Type),
		ServiceID:        string(environment.ServiceID),
		ServiceName:      serviceName,
		SourceEnvName:    sourceEnvName,
		PromoteStatus:    promoteStatus,
		SaaSPortalStatus: saasPortalStatus,
		SaaSPortalURL:    saasPortalURL,
	}

	data, err := json.MarshalIndent(formattedEnvironment, "", "    ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func getPromoteStatus(token, serviceId string, environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) string {
	if !commonutils.CheckIfNilOrEmpty((*string)(environment.SourceEnvironmentID)) {
		promoteRes, err := dataaccess.PromoteServiceEnvironmentStatus(token, serviceId, string(*environment.SourceEnvironmentID))
		if err == nil {
			for _, res := range promoteRes {
				if string(res.TargetEnvironmentID) == string(environment.ID) {
					return res.Status
				}
			}
		}
	}
	return ""
}

func cleanUpCreateFlagsAndArgs(cmd *cobra.Command, args *[]string) {
	// Clean up flags
	_ = cmd.Flags().Set("description", "")
	_ = cmd.Flags().Set("type", "")
	_ = cmd.Flags().Set("source", "")
	_ = cmd.Flags().Set("output", "text")
	_ = cmd.Flags().Set("service-id", "")

	// Clean up arguments by resetting the slice to nil or an empty slice
	*args = nil
}
