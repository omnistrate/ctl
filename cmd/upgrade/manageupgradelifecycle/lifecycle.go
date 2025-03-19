package manageupgradelifecycle

import (
	"fmt"

	"github.com/omnistrate/ctl/cmd/common"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	pauseExample = ` Pausing ongoing upgrade # 
omctl upgrade pause [upgrade-id] `
	resumeExample = ` Resuming paused upgrade # 
omctl upgrade resume [upgrade-id] `
	cancelExample = ` Cancelling uncompleted upgrade # 
omctl upgrade cancel [upgrade-id] `
	notifyCustomerExample = ` Enable customer notifications for a scheduled upgrade # 
omctl upgrade notify-customer [upgrade-id] `
	skipInstancesExample = ` Skip specific instances from an upgrade path #
omctl upgrade skip-instances [upgrade-id] --instances instance-1,instance-2 `
)

var PauseCmd = &cobra.Command{
	Use:          "pause [upgrade-id] [flags]",
	Short:        "Pause an ongoing upgrade",
	Example:      pauseExample,
	RunE:         pause,
	SilenceUsage: true,
}

var ResumeCmd = &cobra.Command{
	Use:          "resume [upgrade-id] [flags]",
	Short:        "Resume a paused upgrade",
	Example:      resumeExample,
	RunE:         resume,
	SilenceUsage: true,
}

var CancelCmd = &cobra.Command{
	Use:          "cancel [upgrade-id] [flags]",
	Short:        "Cancel an uncompleted upgrade",
	Example:      cancelExample,
	RunE:         cancel,
	SilenceUsage: true,
}

var NotifyCustomerCmd = &cobra.Command{
	Use:          "notify-customer [upgrade-id] [flags]",
	Short:        "Enable customer notifications for a scheduled upgrade",
	Example:      notifyCustomerExample,
	RunE:         notifyCustomer,
	SilenceUsage: true,
}

var SkipInstancesCmd = &cobra.Command{
	Use:          "skip-instances [upgrade-id] [flags]",
	Short:        "Skip specific instances from an upgrade path",
	Example:      skipInstancesExample,
	RunE:         skipInstances,
	SilenceUsage: true,
}

func init() {
	PauseCmd.Args = cobra.MinimumNArgs(1)
	ResumeCmd.Args = cobra.MinimumNArgs(1)
	CancelCmd.Args = cobra.MinimumNArgs(1)
	NotifyCustomerCmd.Args = cobra.MinimumNArgs(1)
	SkipInstancesCmd.Args = cobra.MinimumNArgs(1)

	SkipInstancesCmd.Flags().String("instances", "", "Comma-separated list of instance IDs to skip")
	_ = SkipInstancesCmd.MarkFlagRequired("instances")
}

func cancel(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.CancelAction, nil)
}

func pause(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.PauseAction, nil)
}

func resume(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.ResumeAction, nil)
}

func notifyCustomer(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.NotifyCustomerAction, nil)
}

func skipInstances(cmd *cobra.Command, args []string) error {
	instances, err := cmd.Flags().GetString("instances")
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		"resource-ids": utils.SplitCSV(instances),
	}
	return manageLifecycle(cmd, args, model.SkipInstancesAction, payload)
}

func manageLifecycle(cmd *cobra.Command, args []string, action model.UpgradeMaintenanceAction, actionPayload map[string]interface{}) error {
	defer config.CleanupArgsAndFlags(cmd, nil)

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not json
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := fmt.Sprintf("Managing lifecycle of upgrade %s", args[0])
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Search upgrade path
	searchRes, err := dataaccess.SearchInventory(cmd.Context(), token, fmt.Sprintf("upgradepath:%s", args[0]))
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	if searchRes == nil || len(searchRes.UpgradePathResults) == 0 {
		err = fmt.Errorf("%s not found", args[0])
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	found := false
	var serviceID, productTierID string
	for _, upgradePath := range searchRes.UpgradePathResults {
		if upgradePath.Id == args[0] {
			found = true
			serviceID = upgradePath.ServiceId
			productTierID = upgradePath.ProductTierID
			break
		}
	}

	if !found {
		err = fmt.Errorf("%s not found", args[0])
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	upgrade, err := dataaccess.ManageLifecycleWithPayload(cmd.Context(), token, serviceID, productTierID, args[0], action, actionPayload)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	formattedUpgradeStatus := model.UpgradeStatus{
		UpgradeID:  args[0],
		Total:      upgrade.TotalCount,
		Pending:    upgrade.PendingCount,
		InProgress: upgrade.InProgressCount,
		Completed:  upgrade.CompletedCount,
		Failed:     upgrade.FailedCount,
		Scheduled:  utils.FromInt64Ptr(upgrade.ScheduledCount),
		Skipped:    upgrade.SkippedCount,
		Status:     upgrade.Status,
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully managed upgrade lifecycle")

	if output == "json" {
		utils.PrintJSON(formattedUpgradeStatus)
	} else {
		utils.PrintUpgradeStatuses([]*model.UpgradeStatus{&formattedUpgradeStatus})
		fmt.Printf("\nCheck the upgrade status using:\n  omctl upgrade status %s\n", args[0])
	}

	return nil
}
