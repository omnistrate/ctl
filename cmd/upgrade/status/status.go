package status

import (
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/cmd/upgrade/status/detail"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	statusExample = `  # Get upgrade status
  omnistrate-ctl upgrade status <upgrade>`
)

var output string

var Cmd = &cobra.Command{
	Use:          "status",
	Short:        "Get upgrade status",
	Example:      statusExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(detail.Cmd)

	Cmd.Args = cobra.MinimumNArgs(1)

	Cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text|table|json)")
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	res := make([]*model.UpgradeStatus, 0)

	for _, upgradePathID := range args {
		searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("upgradepath:%s", upgradePathID))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if len(searchRes.UpgradePathResults) == 0 {
			err = fmt.Errorf("%s not found", upgradePathID)
			utils.PrintError(err)
			return err
		}

		found := false
		var serviceID, productTierID string
		for _, upgradePath := range searchRes.UpgradePathResults {
			if string(upgradePath.ID) == upgradePathID {
				found = true
				serviceID = string(upgradePath.ServiceID)
				productTierID = string(upgradePath.ProductTierID)
				break
			}
		}

		if !found {
			err = fmt.Errorf("%s not found", upgradePathID)
			utils.PrintError(err)
			return err
		}

		upgrade, err := dataaccess.DescribeUpgradePath(token, serviceID, productTierID, upgradePathID)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		res = append(res, &model.UpgradeStatus{
			UpgradeID:  upgradePathID,
			Total:      upgrade.TotalCount,
			Pending:    upgrade.PendingCount,
			InProgress: upgrade.InProgressCount,
			Completed:  upgrade.CompletedCount,
			Failed:     upgrade.FailedCount,
			Status:     string(upgrade.Status),
		})
	}

	var jsonData []string
	for _, instance := range res {
		data, err := json.MarshalIndent(instance, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		jsonData = append(jsonData, string(data))
	}

	if len(jsonData) == 0 {
		utils.PrintInfo("No upgrades found.")
		return nil
	}

	switch output {
	case "text":
		err = utils.PrintText(jsonData)
		if err != nil {
			return err
		}

		println("\nTo get more details, run the following command(s):")
		for _, r := range res {
			println(fmt.Sprintf("  omnistrate-ctl upgrade status detail %s", r.UpgradeID))
		}
	case "table":
		err = utils.PrintTable(jsonData)
		if err != nil {
			return err
		}

		println("\nTo get more details, run the following command(s):")
		for _, r := range res {
			println(fmt.Sprintf("  omnistrate-ctl upgrade status detail %s", r.UpgradeID))
		}
	case "json":
		fmt.Printf("%+v\n", jsonData)
	default:
		err = fmt.Errorf("invalid output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
