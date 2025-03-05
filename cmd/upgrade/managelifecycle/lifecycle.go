package managelifecycle

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

func init() {
	PauseCmd.Args = cobra.MinimumNArgs(1)
	ResumeCmd.Args = cobra.MinimumNArgs(1)
	CancelCmd.Args = cobra.MinimumNArgs(1)
}
func cancel(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.CancelAction)
}
func pause(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.PauseAction)
}
func resume(cmd *cobra.Command, args []string) error {
	return manageLifecycle(cmd, args, model.ResumeAction)
}

func manageLifecycle(cmd *cobra.Command, args []string, action model.UpgradeMaintenanceAction) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
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
		msg := "Requesting pause action on upgrade..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	formattedUpgradeStatuses := make([]*model.UpgradeStatus, 0)

	for _, upgradePathID := range args {
		searchRes, err := dataaccess.SearchInventory(cmd.Context(), token, fmt.Sprintf("upgradepath:%s", upgradePathID))
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		if len(searchRes.UpgradePathResults) == 0 {
			err = fmt.Errorf("%s not found", upgradePathID)
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		found := false
		var serviceID, productTierID string
		for _, upgradePath := range searchRes.UpgradePathResults {
			if upgradePath.Id == upgradePathID {
				found = true
				serviceID = upgradePath.ServiceId
				productTierID = upgradePath.ProductTierID
				break
			}
		}

		if !found {
			err = fmt.Errorf("%s not found", upgradePathID)
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		upgrade, err := dataaccess.ManageLifecycle(cmd.Context(), token, serviceID, productTierID, upgradePathID, action)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		formattedUpgradeStatuses = append(formattedUpgradeStatuses, &model.UpgradeStatus{
			UpgradeID:  upgradePathID,
			Total:      upgrade.TotalCount,
			Pending:    upgrade.PendingCount,
			InProgress: upgrade.InProgressCount,
			Completed:  upgrade.CompletedCount,
			Failed:     upgrade.FailedCount,
			Scheduled:  utils.FromPtr(upgrade.ScheduledCount),
			Skipped:    upgrade.SkippedCount,
			Status:     upgrade.Status,
		})
	}

	if len(formattedUpgradeStatuses) == 0 {
		utils.HandleSpinnerSuccess(spinner, sm, "No upgrades found")
	} else {
		utils.HandleSpinnerSuccess(spinner, sm, "Upgrade pause request submitted")
	}

	// Print output
	err = utils.PrintTextTableJsonArrayOutput(output, formattedUpgradeStatuses)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if output != "json" {
		println("\nTo get more details, run the following command(s):")
		for _, s := range formattedUpgradeStatuses {
			println(fmt.Sprintf("  omctl upgrade pause detail %s", s.UpgradeID))
		}
	}

	return nil
}
