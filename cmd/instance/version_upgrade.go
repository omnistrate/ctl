package instance

import (
	"errors"
	"fmt"
	"os"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate-oss/omnistrate-ctl/cmd/common"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/config"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/dataaccess"
	"github.com/omnistrate-oss/omnistrate-ctl/internal/utils"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	versionUpgradeExample = `# Issue a version upgrade for an instance with upgrade configuration override to the latest tier version
omctl instance version-upgrade instance-abcd1234 --upgrade-configuration-override /path/to/config.yaml

# Issue a version upgrade to a specific target tier version
omctl instance version-upgrade instance-abcd1234 --upgrade-configuration-override /path/to/config.yaml --target-tier-version v1.2.3

# [HELM ONLY] Use generate-configuration to generate a default deployment instance configuration file based on the current helm values
omctl instance version-upgrade instance-abcd1234 --upgrade-configuration-override existing-config.yaml --generate-configuration

# Example upgrade configuration override YAML file:
# resource-key-1:
#   helmChartValues:
#     key1: value1
#     key2: value2
# resource-key-2:
#   helmChartValues:
#     database:
#       host: new-host
#       port: 5432`
)

var versionUpgradeCmd = &cobra.Command{
	Use:          "version-upgrade [instance-id]",
	Short:        "Issue a version upgrade for a deployment instance",
	Long:         `This command helps you issue a version upgrade for a deployment instance with the specified upgrade configuration override.`,
	Example:      versionUpgradeExample,
	RunE:         runVersionUpgrade,
	SilenceUsage: true,
}

func init() {
	versionUpgradeCmd.Flags().String("upgrade-configuration-override", "", "YAML file containing upgrade configuration override")
	versionUpgradeCmd.Flags().String("target-tier-version", "", "Target tier version for the version upgrade (optional, defaults to latest released tier version)")
	versionUpgradeCmd.Flags().Bool("generate-configuration", false, "Generate a default configuration file based on current helm values")

	versionUpgradeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument (i.e. instance ID)

	var err error
	if err = versionUpgradeCmd.MarkFlagFilename("upgrade-configuration-override"); err != nil {
		return
	}
}

func runVersionUpgrade(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]

	// Retrieve flags
	configOverrideFile, err := cmd.Flags().GetString("upgrade-configuration-override")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Generate configuration override if requested
	generateConfig, err := cmd.Flags().GetBool("generate-configuration")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	targetTierVersion, err := cmd.Flags().GetString("target-tier-version")
	if err != nil {
		utils.PrintError(err)
		return err
	}

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

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		var msg string
		if generateConfig {
			msg = "Generating configuration file"
		} else {
			msg = "Upgrading deployment instance"
			if targetTierVersion != "" {
				msg += fmt.Sprintf(" to target tier version %s", targetTierVersion)
			} else {
				msg += " to latest tier version"
			}
		}
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	if spinner != nil {
		spinner.UpdateMessage("Looking up deployment instance")
	}
	serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// If we have to generate the configuration file, describe the resource instance
	if generateConfig {
		// Update spinner message
		if spinner != nil {
			spinner.UpdateMessage("Processing instance configuration to generate overrides")
		}
		// Describe instance to get current configuration
		instance, err := dataaccess.DescribeResourceInstance(cmd.Context(), token, serviceID, environmentID, instanceID)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		resourceOverrideConfig := make(map[string]openapiclientfleet.ResourceOneOffPatchConfigurationOverride)
		if len(instance.ConsumptionResourceInstanceResult.DetailedNetworkTopology) == 0 {
			utils.HandleSpinnerError(spinner, sm, errors.New("no eligible component topology found for the instance"))
			return errors.New("no eligible component topology found for the instance")
		}

		if spinner != nil {
			spinner.UpdateMessage("Looking up Helm releases for deployment instance to generate overrides")
		}
		for _, resourceVersionSummary := range instance.ResourceVersionSummaries {
			// We only support helm overrides for now
			if resourceVersionSummary.HelmDeploymentConfiguration == nil {
				// Skip resources that are not helm deployments
				continue
			}

			resourceIntfc := instance.ConsumptionResourceInstanceResult.DetailedNetworkTopology[*resourceVersionSummary.ResourceId]

			if resourceIntfc == nil {
				// Skip
				continue
			}

			if resourceMap, ok := resourceIntfc.(map[string]interface{}); ok {
				resourceKey := resourceMap["resourceKey"].(string)
				resourceOverrideConfig[resourceKey] = openapiclientfleet.ResourceOneOffPatchConfigurationOverride{
					HelmChartValues: resourceVersionSummary.HelmDeploymentConfiguration.Values,
				}
			}
		}

		// Write the generated configuration to the specified file
		var marshalledData []byte
		if marshalledData, err = yaml.Marshal(resourceOverrideConfig); err != nil {
			utils.HandleSpinnerError(spinner, sm, errors2.Wrap(err, "failed to marshal generated configuration"))
			return err
		}

		if err = os.WriteFile(configOverrideFile, marshalledData, 0600); err != nil {
			utils.HandleSpinnerError(spinner, sm, errors2.Wrap(err, "failed to write generated configuration to file"))
			return err
		}
		utils.HandleSpinnerSuccess(spinner, sm, fmt.Sprintf("Generated configuration override file: %s", configOverrideFile))
		return nil
	}

	// Read and parse configuration override YAML file
	var resourceOverrideConfig map[string]openapiclientfleet.ResourceOneOffPatchConfigurationOverride
	if configOverrideFile != "" {
		configData, err := os.ReadFile(configOverrideFile)
		if err != nil {
			err = errors2.Wrap(err, "failed to read configuration override file")
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		err = yaml.Unmarshal(configData, &resourceOverrideConfig)
		if err != nil {
			err = errors2.Wrap(err, "failed to parse configuration override YAML")
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}
	}

	// Issue one-off patch
	err = dataaccess.OneOffPatchResourceInstance(cmd.Context(), token,
		serviceID,
		environmentID,
		instanceID,
		resourceOverrideConfig,
		targetTierVersion,
	)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully initiated version upgrade for deployment instance")

	// Search for the instance to get updated details
	searchRes, err := dataaccess.SearchInventory(cmd.Context(), token, fmt.Sprintf("resourceinstance:%s", instanceID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if len(searchRes.ResourceInstanceResults) == 0 {
		err = errors.New("failed to find the upgraded instance")
		utils.PrintError(err)
		return err
	}

	// Format instance
	formattedInstance := formatInstance(&searchRes.ResourceInstanceResults[0], false)

	// Print output
	if err = utils.PrintTextTableJsonOutput(output, formattedInstance); err != nil {
		return err
	}

	return nil
}
