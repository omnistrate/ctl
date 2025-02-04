package upgrade

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd/common"
	"strings"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/upgrade/status"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	upgradeExample = `# Upgrade instances to a specific version
omctl upgrade [instance1] [instance2] --version=2.0

# Upgrade instances to the latest version
omctl upgrade [instance1] [instance2] --version=latest

 # Upgrade instances to the preferred version
omctl upgrade [instance1] [instance2] --version=preferred

# Upgrade instances to a specific version with version name
omctl upgrade [instance1] [instance2] --version-name=v0.1.1`
)

var Cmd = &cobra.Command{
	Use:          "upgrade --version=[version]",
	Short:        "Upgrade Instance Deployments to a newer or older version",
	Long:         `This command helps you upgrade Instance Deployments to a newer or older version.`,
	Example:      upgradeExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	Cmd.AddCommand(status.Cmd)

	Cmd.Args = cobra.MinimumNArgs(1)

	Cmd.Flags().StringP("version", "", "", "Specify the version number to upgrade to. Use 'latest' to upgrade to the latest version. Use 'preferred' to upgrade to the preferred version.")
	Cmd.Flags().StringP("version-name", "", "", "Specify the version name to upgrade to. Use either this flag or the version flag to upgrade to a specific version.")
}

type Args struct {
	ServiceID     string
	ProductTierID string
	SourceVersion string
	TargetVersion string
}

var UpgradePathIDs []string

type Res struct {
	UpgradePathID string
	InstanceIDs   []string
}

func run(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	version, err := cmd.Flags().GetString("version")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	versionName, err := cmd.Flags().GetString("version-name")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate input arguments
	if version == "" && versionName == "" {
		err = errors.New("version or version name is required")
		utils.PrintError(err)
		return err
	}

	if version != "" && versionName != "" {
		err = errors.New("please provide either version or version name, not both")
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
		searchRes, err := dataaccess.SearchInventory(cmd.Context(), token, fmt.Sprintf("resourceinstance:%s", instanceID))
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		if searchRes == nil || len(searchRes.ResourceInstanceResults) == 0 {
			err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID)
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		var found bool
		var serviceID, environmentID, productTierID, sourceVersion, targetVersion string
		for _, instance := range searchRes.ResourceInstanceResults {
			if instance.Id == instanceID {
				serviceID = instance.ServiceId
				environmentID = instance.ServiceEnvironmentId
				productTierID = instance.ProductTierId
				found = true
				break
			}
		}
		if !found {
			err = fmt.Errorf("%s not found. Please check the instance ID and try again", instanceID)
			utils.HandleSpinnerError(spinner, sm, err)
			return nil
		}

		// Find the source version of the instance
		describeRes, err := dataaccess.DescribeResourceInstance(cmd.Context(), token, serviceID, environmentID, instanceID)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
		sourceVersion = describeRes.TierVersion

		// Get the target version
		if version != "" {
			switch version {
			case "latest":
				targetVersion, err = dataaccess.FindLatestVersion(cmd.Context(), token, serviceID, productTierID)
				if err != nil {
					utils.HandleSpinnerError(spinner, sm, err)
					return err
				}
			case "preferred":
				targetVersion, err = dataaccess.FindPreferredVersion(cmd.Context(), token, serviceID, productTierID)
				if err != nil {
					utils.HandleSpinnerError(spinner, sm, err)
					return err
				}
			default:
				targetVersion = version
			}
		} else {
			allVersions, err := dataaccess.ListVersions(cmd.Context(), token, serviceID, productTierID)
			if err != nil {
				utils.HandleSpinnerError(spinner, sm, err)
				return err
			}

			targetVersions := make([]string, 0)
			for _, versionSet := range allVersions.TierVersionSets {
				if versionSet.Name != nil && *versionSet.Name == versionName {
					targetVersions = append(targetVersions, versionSet.Version)
				}
			}

			if len(targetVersions) == 0 {
				err = fmt.Errorf("version name %s not found", versionName)
				utils.HandleSpinnerError(spinner, sm, err)
			}

			if len(targetVersions) > 1 {
				err = fmt.Errorf("multiple versions found for version name %s", versionName)
				utils.HandleSpinnerError(spinner, sm, err)
			}

			targetVersion = targetVersions[0]
		}

		// Check if the target version exists
		_, err = dataaccess.DescribeVersionSet(cmd.Context(), token, serviceID, productTierID, targetVersion)
		if err != nil {
			if strings.Contains(err.Error(), "Version set not found") {
				err = errors.New(fmt.Sprintf("version %s not found", version))
			}
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		// Check if the target is the same as the source
		if sourceVersion == targetVersion {
			err = fmt.Errorf("source version %s is the same as target version for %s", sourceVersion, instanceID)
			utils.HandleSpinnerError(spinner, sm, err)
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
	UpgradePathIDs = make([]string, 0)
	for upgradeArgs, upgradeRes := range upgrades {
		upgradePathID, err := dataaccess.CreateUpgradePath(
			cmd.Context(),
			token,
			upgradeArgs.ServiceID,
			upgradeArgs.ProductTierID,
			upgradeArgs.SourceVersion,
			upgradeArgs.TargetVersion,
			upgradeRes.InstanceIDs,
		)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		upgrades[upgradeArgs].UpgradePathID = string(upgradePathID)
		UpgradePathIDs = append(UpgradePathIDs, string(upgradePathID))
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Upgrade scheduled successfully")

	// Print output
	formattedUpgrades := make([]model.Upgrade, 0)
	for upgradeArgs, upgradeRes := range upgrades {
		formattedUpgrade := model.Upgrade{
			UpgradeID:     upgradeRes.UpgradePathID,
			SourceVersion: upgradeArgs.SourceVersion,
			TargetVersion: upgradeArgs.TargetVersion,
			InstanceIDs:   strings.Join(upgradeRes.InstanceIDs, ","),
		}

		formattedUpgrades = append(formattedUpgrades, formattedUpgrade)
	}

	if output != "json" {
		println("\nThe following upgrades have been scheduled:")
	}

	err = utils.PrintTextTableJsonArrayOutput(output, formattedUpgrades)
	if err != nil {
		return err
	}

	if output != "json" {
		println("\nCheck the upgrade status using the following command(s):")
		for _, upgradeRes := range upgrades {
			fmt.Printf("  omctl upgrade status %s\n", upgradeRes.UpgradePathID)
		}
	}

	return nil
}
