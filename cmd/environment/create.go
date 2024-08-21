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
	// Get flags
	description, _ := cmd.Flags().GetString("description")
	envType, _ := cmd.Flags().GetString("type")
	sourceEnvName, _ := cmd.Flags().GetString("source")
	output, _ := cmd.Flags().GetString("output")
	serviceId, _ := cmd.Flags().GetString("service-id")
	envName := args[len(args)-1]

	if len(args) == 1 && serviceId == "" {
		err := fmt.Errorf("please provide service name or id")
		utils.PrintError(err)
		return err
	}

	if len(args) == 2 && serviceId != "" {
		err := fmt.Errorf("please provide either service name or id, not both")
		utils.PrintError(err)
		return err
	}

	if !slices.Contains([]string{"dev", "qa", "staging", "canary", "prod", "private"}, strings.ToLower(envType)) {
		err := fmt.Errorf("invalid environment type: %s", envType)
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
		msg := "Creating environment..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if the service exists
	searchRes, err := dataaccess.SearchInventory(token, "service:s")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	serviceFound := false
	var serviceName string
	for _, service := range searchRes.ServiceResults {
		if serviceId == service.ID || (len(args) == 2 && strings.EqualFold(service.Name, args[0])) {
			serviceFound = true
			serviceId = service.ID
			serviceName = service.Name
			break
		}
	}

	if !serviceFound {
		err = errors.New("service not found")
		utils.PrintError(err)
		return err
	}

	// Check if source environment exists
	var sourceEnvID string
	if sourceEnvName != "" {
		describeServiceRes, err := dataaccess.DescribeService(token, serviceId)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		sourceEnvFound := false
		for _, env := range describeServiceRes.ServiceEnvironments {
			if strings.EqualFold(env.Name, sourceEnvName) {
				sourceEnvFound = true
				sourceEnvName = env.Name
				sourceEnvID = string(env.ID)
				break
			}
		}

		if !sourceEnvFound {
			err = errors.New("source environment not found. Please provide a valid source environment name")
			utils.PrintError(err)
			return err
		}
	}

	// Create the environment
	if description == "" {
		description = fmt.Sprintf("%s environment for service %s", envType, serviceName)
	}

	visibility := ""
	switch strings.ToLower(envType) {
	case "dev", "qa", "staging", "canary", "private":
		visibility = "PRIVATE"
	case "prod":
		visibility = "PUBLIC"
	default:
		err = fmt.Errorf("invalid environment type: %s", envType)
		utils.PrintError(err)
		return err
	}

	defaultDeploymentConfigID, err := dataaccess.GetDefaultDeploymentConfigID(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	sourceEnvIDPtr := commonutils.ToPtr(sourceEnvID)
	if sourceEnvID == "" {
		sourceEnvIDPtr = nil
	}

	request := serviceenvironmentapi.CreateServiceEnvironmentRequest{
		Name:                    envName,
		Description:             description,
		ServiceID:               serviceenvironmentapi.ServiceID(serviceId),
		Visibility:              serviceenvironmentapi.ServiceVisibility(visibility),
		Type:                    commonutils.ToPtr(serviceenvironmentapi.EnvironmentType(envType)),
		ServiceAuthPublicKey:    commonutils.ToPtr("-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEA2lmruvcEDykT6KbyIJHYCGhCoPUGq+XlCfLWJXlowf4=\n-----END PUBLIC KEY-----"),
		DeploymentConfigID:      serviceenvironmentapi.DeploymentConfigID(defaultDeploymentConfigID),
		AutoApproveSubscription: commonutils.ToPtr(true),
		SourceEnvironmentID:     (*serviceenvironmentapi.ServiceEnvironmentID)(sourceEnvIDPtr),
	}

	environmentID, err := dataaccess.CreateServiceEnvironment(token, request)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if spinner != nil {
		spinner.UpdateMessage("Successfully created environment")
		spinner.Complete()
		sm.Stop()
	}

	// Describe the environment
	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceId, string(environmentID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	saasPortalStatus := ""
	if environment.SaasPortalStatus != nil {
		saasPortalStatus = string(*environment.SaasPortalStatus)
	}
	saasPortalURL := ""
	if environment.SaasPortalURL != nil {
		saasPortalURL = *environment.SaasPortalURL
	}

	// Get promote status
	promoteStatus := ""
	if !commonutils.CheckIfNilOrEmpty((*string)(environment.SourceEnvironmentID)) {
		promoteRes, err := dataaccess.PromoteServiceEnvironmentStatus(token, serviceId, string(*environment.SourceEnvironmentID))
		if err != nil {
			utils.PrintError(err)
			return err
		}
		for _, res := range promoteRes {
			if string(res.TargetEnvironmentID) == string(environment.ID) {
				promoteStatus = res.Status
				break
			}
		}
	}

	// Format the output
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

	var jsonData []string
	data, err := json.MarshalIndent(formattedEnvironment, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	jsonData = append(jsonData, string(data))

	// Print output
	switch output {
	case "text":
		err = utils.PrintText(jsonData)
		if err != nil {
			return err
		}
	case "table":
		err = utils.PrintTable(jsonData)
		if err != nil {
			return err
		}
	case "json":
		_, err = fmt.Fprintf(cmd.OutOrStdout(), "%+v\n", jsonData[0])
		if err != nil {
			utils.PrintError(err)
			return err
		}
	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
