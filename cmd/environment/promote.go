package environment

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	promoteExample = `# Promote environment
omctl environment promote [service-name] [environment-name]

# Promote environment by ID instead of name
omctl environment promote --service-id=[service-id] --environment-id=[environment-id]`
)

var promoteCmd = &cobra.Command{
	Use:          "promote [service-name] [environment-name] [flags]",
	Short:        "Promote a environment",
	Long:         `This command helps you promote a environment in your service.`,
	Example:      promoteExample,
	RunE:         runPromote,
	SilenceUsage: true,
}

func init() {
	promoteCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	promoteCmd.Flags().StringP("environment-id", "", "", "Environment ID. Required if environment name is not provided")
}

func runPromote(cmd *cobra.Command, args []string) error {
	defer utils.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	output, _ := cmd.Flags().GetString("output")
	serviceID, _ := cmd.Flags().GetString("service-id")
	environmentID, _ := cmd.Flags().GetString("environment-id")

	// Validate input arguments
	if err := validatePromoteArguments(args, serviceID, environmentID); err != nil {
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
		spinner = sm.AddSpinner("Promoting environment...")
		sm.Start()
	}

	// Check if the environment exists
	serviceID, serviceName, environmentID, _, err = getServiceEnvironment(token, serviceID, serviceName, environmentID, environmentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Promote the environment
	if err = dataaccess.PromoteServiceEnvironment(token, serviceID, environmentID); err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe the promoted environment
	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceID, environmentID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Get promote status and format output
	formattedPromotions, err := formatPromoteStatus(token, serviceID, environmentID, serviceName, environment)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Handle output based on format
	utils.HandleSpinnerSuccess(spinner, sm, "Successfully promoted environment")

	if err = utils.PrintTextTableJsonArrayOutput(output, formattedPromotions); err != nil {
		return err
	}

	return nil
}

// Helper functions

func validatePromoteArguments(args []string, serviceID, environmentID string) error {
	if len(args) == 0 && (serviceID == "" || environmentID == "") {
		return fmt.Errorf("please provide the service name and environment name or the service ID and environment ID")
	}
	if len(args) > 0 && len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [environment-name]", strings.Join(args, " "))
	}
	return nil
}

func getServiceEnvironment(token, serviceIDArg, serviceNameArg, environmentIDArg, environmentNameArg string) (serviceID, serviceName, environmentID, environmentName string, err error) {
	services, err := dataaccess.ListServices(token)
	if err != nil {
		return
	}

	serviceEnvironmentsMap := make(map[string]map[string]bool)

	for _, service := range services.Services {
		if string(service.ID) != serviceIDArg && !strings.EqualFold(service.Name, serviceNameArg) {
			continue
		}

		for _, environment := range service.ServiceEnvironments {
			if string(environment.ID) != environmentIDArg && !strings.EqualFold(environment.Name, environmentNameArg) {
				continue
			}

			if _, exists := serviceEnvironmentsMap[string(service.ID)]; !exists {
				serviceEnvironmentsMap[string(service.ID)] = make(map[string]bool)
			}
			serviceEnvironmentsMap[string(service.ID)][string(environment.ID)] = true

			serviceID = string(service.ID)
			environmentID = string(environment.ID)
			serviceName = service.Name
			environmentName = environment.Name
		}
	}

	if len(serviceEnvironmentsMap) == 0 {
		err = errors.New("environment not found. Please check the input values and try again")
		return
	}
	if len(serviceEnvironmentsMap) > 1 || len(serviceEnvironmentsMap[serviceID]) > 1 {
		err = errors.New("multiple environments found. Please provide the service ID and environment ID instead of the names")
		return
	}
	return
}

func formatPromoteStatus(token, serviceID, environmentID, serviceName string, environment *serviceenvironmentapi.DescribeServiceEnvironmentResult) ([]model.Promotion, error) {
	promotions, err := dataaccess.PromoteServiceEnvironmentStatus(token, serviceID, environmentID)
	if err != nil {
		return nil, err
	}

	var formattedPromotions []model.Promotion
	for _, promotion := range promotions {
		targetEnvID := string(promotion.TargetEnvironmentID)
		targetEnv, err := dataaccess.DescribeServiceEnvironment(token, serviceID, targetEnvID)
		if err != nil {
			return nil, err
		}

		formattedPromotion := model.Promotion{
			ServiceID:             string(environment.ServiceID),
			ServiceName:           serviceName,
			SourceEnvironmentID:   string(environment.ID),
			SourceEnvironmentName: environment.Name,
			TargetEnvID:           targetEnvID,
			TargetEnvName:         targetEnv.Name,
			PromoteStatus:         promotion.Status,
		}
		formattedPromotions = append(formattedPromotions, formattedPromotion)
	}

	if len(formattedPromotions) == 0 {
		return nil, fmt.Errorf("source environment %s is not linked to any target environments", environment.Name)
	}

	return formattedPromotions, nil
}
