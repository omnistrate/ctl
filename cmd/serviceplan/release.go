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
	releaseExample = `# Release service plan
omnistrate service-plan release [service-name] [plan-name]

# Release service plan by ID instead of name
omnistrate service-plan release --service-id [plan-id] --plan-id [plan-id]`
)

var releaseCmd = &cobra.Command{
	Use:          "release [service-name] [plan-name] [flags]",
	Short:        "Release a service plan",
	Long:         `This command helps you release a service plan from your service.`,
	Example:      releaseExample,
	RunE:         runRelease,
	SilenceUsage: true,
}

func init() {
	releaseCmd.Flags().String("release-description", "", "Set custom release description for this release version")
	releaseCmd.Flags().Bool("release-as-preferred", false, "Release the service plan as preferred")
	releaseCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	releaseCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	releaseCmd.Flags().StringP("plan-id", "", "", "Plan ID. Required if plan name is not provided")
}

func runRelease(cmd *cobra.Command, args []string) error {
	// Get flags
	releaseDescription, _ := cmd.Flags().GetString("release-description")
	releaseAsPreferred, _ := cmd.Flags().GetBool("release-as-preferred")
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
		msg := "Releasing service plan..."
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

	// Find service api id
	productTier, err := dataaccess.DescribeProductTier(token, serviceId, planId)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	serviceModel, err := dataaccess.DescribeServiceModel(token, serviceId, string(productTier.ServiceModelID))
	if err != nil {
		utils.PrintError(err)
		return err
	}
	serviceApiId := string(serviceModel.ServiceAPIID)

	// Release service plan
	var releaseDescriptionPtr *string
	if releaseDescription != "" {
		releaseDescriptionPtr = &releaseDescription
	}
	err = dataaccess.ReleaseServicePlan(token, serviceId, serviceApiId, planId, releaseDescriptionPtr, releaseAsPreferred)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if output != "json" {
		spinner.UpdateMessage("Successfully released service plan")
		spinner.Complete()
		sm.Stop()
	}

	// Search it in the inventory
	searchRes, err = dataaccess.SearchInventory(token, fmt.Sprintf("serviceplan:%s", planId))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	latestVersion, err := dataaccess.FindLatestVersion(token, serviceId, planId)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var formattedServicePlan model.ServicePlan
	for _, servicePlan := range searchRes.ServicePlanResults {
		if string(servicePlan.ServiceID) == serviceId && servicePlan.ID == planId && servicePlan.Version == latestVersion {
			envType := ""
			if servicePlan.ServiceEnvironmentType != nil {
				envType = string(*servicePlan.ServiceEnvironmentType)
			}
			versionName := ""
			if servicePlan.VersionName != nil {
				versionName = *servicePlan.VersionName
			}
			formattedServicePlan = model.ServicePlan{
				PlanID:             servicePlan.ID,
				PlanName:           servicePlan.Name,
				ServiceID:          string(servicePlan.ServiceID),
				ServiceName:        servicePlan.ServiceName,
				Environment:        envType,
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
