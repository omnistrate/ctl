package environment

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	promoteExample = `# Promote environment
omnistrate environment promote [service-name] [environment-name]

# Promote environment by ID instead of name
omnistrate environment promote --service-id [service-id] --environment-id [environment-id]`
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
	promoteCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	promoteCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	promoteCmd.Flags().StringP("environment-id", "", "", "Environment ID. Required if environment name is not provided")
}

func runPromote(cmd *cobra.Command, args []string) error {
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
		msg := "Promoting environment..."
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

	// Promote the environment
	err = dataaccess.PromoteServiceEnvironment(token, serviceId, environmentId)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if spinner != nil {
		spinner.UpdateMessage("Successfully promoted environment")
		spinner.Complete()
		sm.Stop()
	}

	// Describe the environment
	environment, err := dataaccess.DescribeServiceEnvironment(token, serviceId, string(environmentId))
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
