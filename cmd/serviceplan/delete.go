package serviceplan

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	deleteExample = `# Delete service plan
omnistrate service-plan delete [service-name] [plan-name]

# Delete service plan by ID instead of name
omnistrate service-plan delete --service-id [service-id] --plan-id [plan-id]`
)

var deleteCmd = &cobra.Command{
	Use:          "delete [service-name] [plan-name] [flags]",
	Short:        "Delete a service plan",
	Long:         `This command helps you delete a service plan from your service.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
	deleteCmd.Flags().StringP("service-id", "", "", "Service ID. Required if service name is not provided")
	deleteCmd.Flags().StringP("plan-id", "", "", "Plan ID. Required if plan name is not provided")
}

func runDelete(cmd *cobra.Command, args []string) error {
	// Get flags
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
		msg := "Deleting service plan..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if the service plan exists
	searchRes, err := dataaccess.SearchInventory(token, "serviceplan:pt")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	servicePlansMap := make(map[string]map[string]bool)
	for _, plan := range searchRes.ServicePlanResults {
		if (string(plan.ServiceID) == serviceId || (len(args) == 2 && strings.EqualFold(plan.ServiceName, args[0]))) &&
			(plan.ID == planId || (len(args) == 2 && strings.EqualFold(plan.Name, args[1]))) {
			if _, ok := servicePlansMap[string(plan.ServiceID)]; !ok {
				servicePlansMap[string(plan.ServiceID)] = make(map[string]bool)
			}
			servicePlansMap[string(plan.ServiceID)][plan.ID] = true
			serviceId = string(plan.ServiceID)
			planId = plan.ID
		}
	}
	if len(servicePlansMap) == 0 {
		err = errors.New("service plan not found. Please check the input values and try again")
		utils.PrintError(err)
		return err
	}
	if len(servicePlansMap) > 1 || len(servicePlansMap[serviceId]) > 1 {
		err = errors.New("multiple service plans found. Please provide the service ID and plan ID instead of the names")
		utils.PrintError(err)
		return err
	}

	err = dataaccess.DeleteServicePlan(token, serviceId, planId)
	if err != nil {
		spinner.Error()
		sm.Stop()
		utils.PrintError(err)
		return err
	}

	if output != "json" {
		spinner.UpdateMessage("Successfully deleted service plan")
		spinner.Complete()
		sm.Stop()
	}

	return nil
}
