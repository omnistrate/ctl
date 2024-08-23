package detail

import (
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/model"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	detailExample = `  # Get upgrade status detail
  omctl upgrade status detail <upgrade>`
)

var output string

var Cmd = &cobra.Command{
	Use:          "detail",
	Short:        "Get upgrade status detail",
	Example:      detailExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.Args = cobra.ExactArgs(1)

	Cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text|table|json)")
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	res := make([]*model.UpgradeStatusDetail, 0)

	upgradePathID := args[0]
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

	instanceUpgrades, err := dataaccess.ListEligibleInstancesPerUpgrade(token, serviceID, productTierID, upgradePathID)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	for _, instanceUpgrade := range instanceUpgrades {
		startTime := ""
		if instanceUpgrade.UpgradeStartTime != nil {
			startTime = *instanceUpgrade.UpgradeStartTime
		}

		endTime := ""
		if instanceUpgrade.UpgradeEndTime != nil {
			endTime = *instanceUpgrade.UpgradeEndTime
		}
		res = append(res, &model.UpgradeStatusDetail{
			UpgradeID:        upgradePathID,
			InstanceID:       string(instanceUpgrade.InstanceID),
			UpgradeStatus:    string(instanceUpgrade.Status),
			UpgradeStartTime: startTime,
			UpgradeEndTime:   endTime,
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
		fmt.Printf("%+v\n", jsonData)
	default:
		err = fmt.Errorf("invalid output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
