package environment

import (
	"encoding/json"
	"fmt"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	describeExample = `# Describe environment
omnistrate environment describe [service-name] [environment-name]

# Describe environment by ID instead of name
omnistrate environment describe --service-id [service-id] --environment-id [environment-id]`
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
	// Get flags
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

	// Check if the environment exists
	services, err := dataaccess.ListServices(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var serviceName string
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

	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceId, environmentId)
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

	// Get source environment name
	sourceEnvName := ""
	if !commonutils.CheckIfNilOrEmpty((*string)(environment.SourceEnvironmentID)) {
		sourceEnv, err := dataaccess.DescribeServiceEnvironment(token, serviceId, string(*environment.SourceEnvironmentID))
		if err != nil {
			utils.PrintError(err)
			return err
		}
		sourceEnvName = sourceEnv.Name
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

	data, err := json.MarshalIndent(formattedEnvironment, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
