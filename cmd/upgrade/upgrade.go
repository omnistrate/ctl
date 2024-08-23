package upgrade

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/upgrade/status"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

const (
	upgradeExample = `  # Upgrade instances to a specific version
  omctl upgrade <instance1> <instance2> --version 2.0

  # Upgrade instances to the latest version
  omctl upgrade <instance1> <instance2> --version latest

 # Upgrade instances to the preferred version
  omctl upgrade <instance1> <instance2> --version preferred`
)

var version string
var output string

var Cmd = &cobra.Command{
	Use:          "upgrade [--version VERSION]",
	Short:        "Upgrade instance to a newer or older version",
	Long:         `This command helps you upgrade instances to a newer or older version.`,
	Example:      upgradeExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(status.Cmd)

	Cmd.Args = cobra.MinimumNArgs(1)

	Cmd.Flags().StringVarP(&version, "version", "", "", "Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version. Use 'preferred' to upgrade to the preferred version.")
	Cmd.Flags().StringVarP(&output, "output", "o", "text", "Output format (text|table|json)")

	err := Cmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

type Args struct {
	ServiceID     string
	ProductTierID string
	SourceVersion string
	TargetVersion string
}

type Res struct {
	UpgradePathID string
	InstanceIDs   []string
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Scheduling upgrade for all instances"
		if len(args) == 1 {
			msg = fmt.Sprintf("Scheduling upgrade for %s", args[0])
		}
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	upgrades := make(map[Args]*Res)
	for _, instanceID := range args {
		// Check if the instance exists
		searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", instanceID))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if searchRes == nil || len(searchRes.ResourceInstanceResults) == 0 {
			err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID)
			utils.PrintError(err)
			return err
		}

		var found bool
		var serviceID, environmentID, productTierID, sourceVersion, targetVersion string
		for _, instance := range searchRes.ResourceInstanceResults {
			if instance.ID == instanceID {
				serviceID = string(instance.ServiceID)
				environmentID = string(instance.ServiceEnvironmentID)
				productTierID = string(instance.ProductTierID)
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID)
			utils.PrintError(err)
			return nil
		}

		// Find the source version of the instance
		describeRes, err := dataaccess.DescribeResourceInstance(token, serviceID, environmentID, instanceID)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		sourceVersion = describeRes.TierVersion

		// Get the target version
		switch version {
		case "latest":
			targetVersion, err = dataaccess.FindLatestVersion(token, serviceID, productTierID)
			if err != nil {
				utils.PrintError(err)
				return err
			}
		case "preferred":
			targetVersion, err = dataaccess.FindPreferredVersion(token, serviceID, productTierID)
			if err != nil {
				utils.PrintError(err)
				return err
			}
		default:
			targetVersion = version
		}

		// Check if the target version exists
		_, err = dataaccess.DescribeVersionSet(token, serviceID, productTierID, targetVersion)
		if err != nil {
			if strings.Contains(err.Error(), "Version set not found") {
				err = errors.New(fmt.Sprintf("version %s not found", version))
			}
			utils.PrintError(err)
			return err
		}

		// Check if the target is the same as the source
		if sourceVersion == targetVersion {
			err = fmt.Errorf("source version %s is the same as target version for %s", sourceVersion, instanceID)
			utils.PrintError(err)
			return err
		}

		if upgrades[Args{
			ServiceID:     serviceID,
			ProductTierID: productTierID,
			SourceVersion: sourceVersion,
			TargetVersion: targetVersion,
		}] == nil {
			upgrades[Args{
				ServiceID:     serviceID,
				ProductTierID: productTierID,
				SourceVersion: sourceVersion,
				TargetVersion: targetVersion,
			}] = &Res{
				InstanceIDs: make([]string, 0),
			}
		}

		upgrades[Args{
			ServiceID:     serviceID,
			ProductTierID: productTierID,
			SourceVersion: sourceVersion,
			TargetVersion: targetVersion,
		}].InstanceIDs = append(upgrades[Args{
			ServiceID:     serviceID,
			ProductTierID: productTierID,
			SourceVersion: sourceVersion,
			TargetVersion: targetVersion,
		}].InstanceIDs, instanceID)
	}

	// Create upgrade path
	for upgradeArgs, upgradeRes := range upgrades {
		upgradePathID, err := dataaccess.CreateUpgradePath(token, upgradeArgs.ServiceID, upgradeArgs.ProductTierID, upgradeArgs.SourceVersion, upgradeArgs.TargetVersion, upgradeRes.InstanceIDs)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		upgrades[upgradeArgs].UpgradePathID = string(upgradePathID)
	}

	if spinner != nil {
		spinner.Complete()
		sm.Stop()
	}

	// Print output
	switch output {
	case "text":
		println("\nThe following upgrades have been scheduled:")

		printTable(upgrades)

		println("\nCheck the upgrade status using the following command(s):")
		for _, upgradeRes := range upgrades {
			fmt.Printf("  omctl upgrade status %s\n", upgradeRes.UpgradePathID)
		}
	case "json":
		printJSON(upgrades)
	default:
		err = fmt.Errorf("invalid output format %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}

func printTable(upgrades map[Args]*Res) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	_, err := fmt.Fprintln(w, "Upgrade ID\tSource Version\tTarget Version\tInstance IDs")
	if err != nil {
		return
	}

	for args, res := range upgrades {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			res.UpgradePathID,
			args.SourceVersion,
			args.TargetVersion,
			strings.Join(res.InstanceIDs, ", "),
		)
		if err != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}

func printJSON(upgrades map[Args]*Res) {
	slices := make([]interface{}, 0)
	for args, res := range upgrades {
		slices = append(slices, map[string]interface{}{
			"upgrade_id":     res.UpgradePathID,
			"source_version": args.SourceVersion,
			"target_version": args.TargetVersion,
			"instance_ids":   res.InstanceIDs,
		})
	}

	utils.PrintJSON(slices)
}
