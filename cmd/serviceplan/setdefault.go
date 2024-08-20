package serviceplan

import (
	"encoding/json"
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"strings"
)

const (
	setDefaultExample = `# Set service plan as default
omnistrate service-plan set-default [service-name] [plan-name] --version [version]

# Set  service plan as default by ID instead of name
omnistrate service-plan set-default --service-id [plan-id] --plan-id [plan-id] --version [version]`
)

var setDefaultCmd = &cobra.Command{
	Use:   "set-default [service-name] [plan-name] [--version=VERSION] [flags]",
	Short: "Set a service plan as default",
	Long: `This command helps you set a service plan as default for your service.
By setting a service plan as default, you can ensure that new instances of the service are created with the default plan.`,
	Example:      setDefaultExample,
	RunE:         runSetDefault,
	SilenceUsage: true,
}

func init() {
	setDefaultCmd.Flags().String("version", "", "Specify the version number to set the default to. Use 'latest' to set the latest version as default.")
	setDefaultCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	setDefaultCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	setDefaultCmd.Flags().StringP("plan-id", "", "", "Plan ID. Required if plan name is not provided")
}

func runSetDefault(cmd *cobra.Command, args []string) error {
	// Get flags
	version, _ := cmd.Flags().GetString("version")
	output, _ := cmd.Flags().GetString("output")
	serviceId, _ := cmd.Flags().GetString("service-id")
	planId, _ := cmd.Flags().GetString("plan-id")

	if len(args) == 0 {
		// Check if service ID and plan ID are provided
		if serviceId == "" || planId == "" {
			err := fmt.Errorf("please provide the service name and plan name or the service ID and plan ID")
			utils.PrintError(err)
			return err
		}
	}

	if len(args) > 0 && len(args) != 2 {
		err := fmt.Errorf("invalid arguments: %s. Need 2 arguments: [service-name] [plan-name]", strings.Join(args, " "))
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
		msg := "Setting default service plan..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Get service ID and plan ID
	searchRes, err := dataaccess.SearchInventory(token, "service:s")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	serviceFound := false
	for _, service := range searchRes.ServiceResults {
		if (len(args) > 0 && strings.EqualFold(service.Name, args[0])) || service.ID == serviceId {
			serviceId = service.ID
			serviceFound = true
		}
	}

	if !serviceFound {
		err = fmt.Errorf("service not found. Please check input values and try again")
		utils.PrintError(err)
		return err
	}

	servicePlanFound := false
	describeServiceRes, err := dataaccess.DescribeService(serviceId, token)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	for _, env := range describeServiceRes.ServiceEnvironments {
		for _, servicePlan := range env.ServicePlans {
			if (len(args) > 0 && strings.EqualFold(servicePlan.Name, args[1])) || string(servicePlan.ProductTierID) == planId {
				planId = string(servicePlan.ProductTierID)
				servicePlanFound = true
			}
		}
	}

	if !servicePlanFound {
		err = fmt.Errorf("service plan not found. Please check input values and try again")
		utils.PrintError(err)
		return err
	}

	// Get the target version
	var targetVersion string
	switch version {
	case "latest":
		targetVersion, err = dataaccess.FindLatestVersion(token, serviceId, planId)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	default:
		targetVersion = version
	}

	// Release service plan
	_, err = dataaccess.SetDefaultServicePlan(token, serviceId, planId, targetVersion)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if output != "json" {
		spinner.UpdateMessage("Successfully setDefaultd service plan")
		spinner.Complete()
		sm.Stop()
	}

	// Search it in the inventory
	searchRes, err = dataaccess.SearchInventory(token, fmt.Sprintf("serviceplan:%s", planId))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var formattedServicePlan model.ServicePlan
	for _, servicePlan := range searchRes.ServicePlanResults {
		if string(servicePlan.ServiceID) == serviceId && servicePlan.ID == planId && servicePlan.Version == targetVersion {
			versionName := ""
			if servicePlan.VersionName != nil {
				versionName = *servicePlan.VersionName
			}
			formattedServicePlan = model.ServicePlan{
				PlanID:             servicePlan.ID,
				PlanName:           servicePlan.Name,
				ServiceID:          string(servicePlan.ServiceID),
				ServiceName:        servicePlan.ServiceName,
				Environment:        servicePlan.ServiceEnvironmentName,
				Version:            servicePlan.Version,
				ReleaseDescription: versionName,
				VersionSetStatus:   servicePlan.VersionSetStatus,
				DeploymentType:     string(servicePlan.DeploymentType),
				TenancyType:        string(servicePlan.TenancyType),
			}
		}
	}

	var jsonData []string
	data, err := json.MarshalIndent(formattedServicePlan, "", "    ")
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
