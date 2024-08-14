package detail

import (
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

const (
	detailLong = ``

	detailExample = `  # Get upgrade status detail
  omnistrate-ctl upgrade status detail <upgrade>`
)

var output string

var Cmd = &cobra.Command{
	Use:          "detail",
	Short:        "Get upgrade status detail",
	Long:         detailLong,
	Example:      detailExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.Args = cobra.ExactArgs(1)

	Cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format. One of: text, json")
}

type Res struct {
	UpgradeID        string
	InstanceID       string
	UpgradeStartTime string
	UpgradeEndTime   string
	UpgradeStatus    string
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	res := make([]*Res, 0)

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
		res = append(res, &Res{
			UpgradeID:        upgradePathID,
			InstanceID:       string(instanceUpgrade.InstanceID),
			UpgradeStatus:    string(instanceUpgrade.Status),
			UpgradeStartTime: startTime,
			UpgradeEndTime:   endTime,
		})
	}

	switch output {
	case "text":
		printTable(res)
	case "json":
		utils.PrintJSON(res)
	default:
		err = fmt.Errorf("invalid output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}

func printTable(res []*Res) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Instance ID\tStatus\tStart Time\tEnd Time")

	for _, r := range res {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			r.InstanceID,
			r.UpgradeStatus,
			r.UpgradeStartTime,
			r.UpgradeEndTime,
		)
	}

	err := w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
