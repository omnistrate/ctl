package status

import (
	"fmt"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

const (
	statusLong = ``

	statusExample = `  # Get upgrade status
  omnistrate-ctl upgrade status <upgrade1> <upgrade2>`
)

var output string

var Cmd = &cobra.Command{
	Use:          "status <upgrade>",
	Short:        "Get upgrade status",
	Long:         statusLong,
	Example:      statusExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.Args = cobra.MinimumNArgs(1)

	Cmd.Flags().StringVarP(&output, "output", "o", "table", "Output format. One of: table, json")
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

	allUpgrades, err := dataaccess.ListAllUpgradePaths(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	allUpgradesMap := make(map[string]*inventoryapi.UpgradePath)
	for _, upgradePath := range allUpgrades {
		allUpgradesMap[string(upgradePath.UpgradePathID)] = upgradePath
	}

	res := make([]*Res, 0)

	for _, upgradePathID := range args {
		upgradePath, ok := allUpgradesMap[upgradePathID]
		if !ok {
			err = errors.New(fmt.Sprintf("%s not found", upgradePathID))
			utils.PrintError(err)
			return err
		}

		instanceUpgrades, err := dataaccess.ListEligibleInstancesPerUpgrade(token, string(upgradePath.ServiceID), string(upgradePath.ProductTierID), string(upgradePath.UpgradePathID))
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
				UpgradeID:        string(upgradePath.UpgradePathID),
				InstanceID:       string(instanceUpgrade.InstanceID),
				UpgradeStartTime: startTime,
				UpgradeEndTime:   endTime,
				UpgradeStatus:    string(instanceUpgrade.Status),
			})
		}
	}

	switch output {
	case "table":
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

	fmt.Fprintln(w, "Upgrade ID\tInstance ID\tStart Time\tEnd Time\tStatus")

	for _, r := range res {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			r.UpgradeID,
			r.InstanceID,
			r.UpgradeStartTime,
			r.UpgradeEndTime,
			r.UpgradeStatus,
		)
	}

	err := w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
