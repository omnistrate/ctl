package status

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd/upgrade/status/detail"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

const (
	statusLong = ``

	statusExample = `  # Get upgrade status
  omnistrate-ctl upgrade status <upgrade>`
)

var output string

var Cmd = &cobra.Command{
	Use:          "status <upgrade>",
	Short:        "Get upgrade status",
	Long:         statusLong,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(detail.Cmd)

	Cmd.Example = getExample()

	Cmd.Args = cobra.MinimumNArgs(1)

	Cmd.Flags().StringVarP(&output, "output", "o", "table", "Output format. One of: table, json")
}

func getExample() (example string) {
	example += statusExample + "\n\n"
	for _, cmd := range Cmd.Commands() {
		example += cmd.Example + "\n\n"
	}
	return example
}

type Res struct {
	UpgradeID  string
	Total      int
	Pending    int
	InProgress int
	Completed  int
	Failed     int
	Status     string
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	res := make([]*Res, 0)

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

		res = append(res, &Res{
			UpgradeID:  upgradePathID,
			Total:      upgrade.TotalCount,
			Pending:    upgrade.PendingCount,
			InProgress: upgrade.InProgressCount,
			Completed:  upgrade.CompletedCount,
			Failed:     upgrade.FailedCount,
			Status:     string(upgrade.Status),
		})
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

	println("\nTo get more details, run the following command(s):")
	for _, r := range res {
		println(fmt.Sprintf("  omnistrate-ctl upgrade status detail %s", r.UpgradeID))
	}

	return nil
}

func printTable(res []*Res) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "Upgrade ID\tTotal\tPending\tIn Progress\tCompleted\tFailed\tStatus")

	for _, r := range res {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\t%d\t%s\n",
			r.UpgradeID,
			r.Total,
			r.Pending,
			r.InProgress,
			r.Completed,
			r.Failed,
			r.Status,
		)
	}

	err := w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
