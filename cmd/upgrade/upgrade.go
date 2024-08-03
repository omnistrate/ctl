package upgrade

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"sync"
)

const (
	upgradeLong = ``

	upgradeExample = `  # Upgrade instances to a specific version
  omnistrate-ctl upgrade <instance1> <instance2> --version 1.2.3

  # Upgrade instances to the latest version
  omnistrate-ctl upgrade <instance1> <instance2> --version latest`
)

var version string

var UpgradeCmd = &cobra.Command{
	Use:          "upgrade <instance> [--version VERSION]",
	Short:        "Upgrade instance to a newer version or an older version.",
	Long:         upgradeLong,
	Example:      upgradeExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	UpgradeCmd.Args = cobra.MinimumNArgs(1)

	UpgradeCmd.Flags().StringVarP(&version, "version", "", "", "Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version.")

	err := UpgradeCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

type Upgrade struct {
	InstanceID    string
	ServiceID     string
	EnvironmentID string
	ProductTierID string
	SourceVersion string
	TargetVersion string
	UpgradePathID string
	Spinner       *ysmrr.Spinner
}

func run(cmd *cobra.Command, args []string) error {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	sm := ysmrr.NewSpinnerManager()

	upgrades := make(map[string]*Upgrade)
	for _, instanceID := range args {
		upgrades[instanceID].Spinner = sm.AddSpinner(fmt.Sprintf("preparing %s", instanceID))

		// Check if the instance exists
		searchRes, err := dataaccess.SearchInventory(token, fmt.Sprintf("resourceinstance:%s", instanceID))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if searchRes == nil || len(searchRes.ResourceInstanceResults) == 0 {
			utils.PrintError(fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID))
			return nil
		}

		var found bool
		for _, instance := range searchRes.ResourceInstanceResults {
			if instance.ID == instanceID {
				upgrades[instanceID].ServiceID = string(instance.ServiceID)
				upgrades[instanceID].EnvironmentID = string(instance.ServiceEnvironmentID)
				upgrades[instanceID].ProductTierID = string(instance.ProductTierID)
				break
			}
		}
		if !found {
			utils.PrintError(fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID))
			return nil
		}

		// Find the source version of the instance
		describeRes, err := dataaccess.DescribeResourceInstance(token, upgrades[instanceID].ServiceID, upgrades[instanceID].EnvironmentID, instanceID)
		if err != nil {
			utils.PrintError(err)
			return err
		}
		upgrades[instanceID].SourceVersion = describeRes.TierVersion

		// Check if the target version exists
		if version == "latest" {
			upgrades[instanceID].TargetVersion, err = dataaccess.FindLatestVersion(token, upgrades[instanceID].ServiceID, upgrades[instanceID].ProductTierID)
			if err != nil {
				utils.PrintError(err)
				return err
			}
		} else {
			upgrades[instanceID].TargetVersion = version
		}

		// Check if the target version exists
		_, err = dataaccess.DescribeVersionSet(token, upgrades[instanceID].ServiceID, upgrades[instanceID].ProductTierID, upgrades[instanceID].TargetVersion)
		if err != nil {
			utils.PrintError(err)
			return err
		}
	}
	sm.Start()

	// Create upgrade path
	var wg sync.WaitGroup
	for _, instanceID := range args {
		upgradePathID, err := dataaccess.CreateUpgradePath(token, upgrades[instanceID].ServiceID, upgrades[instanceID].ProductTierID, upgrades[instanceID].SourceVersion, upgrades[instanceID].TargetVersion, instanceID)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("upgrading %s", instanceID))
		upgrades[instanceID].UpgradePathID = string(upgradePathID)

		wg.Add(1)
	}

	// Check if upgrade completed
	for _, instanceID := range args {
		go func(instanceID string) {
			for {
				upgradePath, err := dataaccess.DescribeUpgradePath(token, upgrades[instanceID].ServiceID, upgrades[instanceID].ProductTierID, upgrades[instanceID].UpgradePathID)
				if err != nil {
					utils.PrintError(err)
					return
				}

				switch upgradePath.Status {
				case "PENDING":
					upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("%s pending", instanceID))
				case "IN_PROGRESS":
					upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("%s in progress", instanceID))
				case "COMPLETE":
					upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("%s completed", instanceID))
					upgrades[instanceID].Spinner.Complete()
					wg.Done()
					break
				case "FAILED":
					upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("%s failed", instanceID))
					upgrades[instanceID].Spinner.Error()
					wg.Done()
					break
				case "PAUSED":
					upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("%s paused", instanceID))
					upgrades[instanceID].Spinner.Error()
					wg.Done()
					break
				case "CANCELLED":
					upgrades[instanceID].Spinner.UpdateMessage(fmt.Sprintf("%s cancelled", instanceID))
					upgrades[instanceID].Spinner.Error()
					wg.Done()
					break
				default:
					err := fmt.Errorf("unknown status: %s", upgradePath.Status)
					utils.PrintError(err)
					wg.Done()
					break
				}
			}

		}(instanceID)
	}

	wg.Wait()
	sm.Stop()

	return nil
}
