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
	patchExample = `# Issue a one-off patch to an instance with configuration override
omctl instance patch instance-abcd1234 --configuration-override /path/to/config.yaml

# Issue a one-off patch with target tier version
omctl instance patch instance-abcd1234 --configuration-override /path/to/config.yaml --target-tier-version v1.2.3

# Example configuration override YAML file:
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

var patchCmd = &cobra.Command{
	Use:          "patch [instance-id]",
	Short:        "Issue a one-off patch to an instance",
	Long:         `This command helps you issue a one-off patch to a resource instance with configuration overrides.`,
	Example:      patchExample,
	RunE:         runPatch,
	SilenceUsage: true,
}

func init() {
	patchCmd.Flags().String("configuration-override", "", "YAML file containing resource configuration overrides")
	patchCmd.Flags().String("target-tier-version", "", "Target tier version for the patch")
	patchCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	patchCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	var err error
	if err = patchCmd.MarkFlagRequired("configuration-override"); err != nil {
		return
	}
	if err = patchCmd.MarkFlagFilename("configuration-override"); err != nil {
		return
	}
}

func runPatch(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]

	// Retrieve flags
	configOverrideFile, err := cmd.Flags().GetString("configuration-override")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if configOverrideFile == "" {
		err = errors.New("configuration-override is required")
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

	// Validate output flag
	if output != "json" {
		err = errors.New("only json output is supported")
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
		msg := "Applying one-off patch..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Check if instance exists
	serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Read and parse configuration override YAML file
	configData, err := os.ReadFile(configOverrideFile)
	if err != nil {
		err = errors2.Wrap(err, "failed to read configuration override file")
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	var resourceOverrideConfig map[string]openapiclientfleet.ResourceOneOffPatchConfigurationOverride
	err = yaml.Unmarshal(configData, &resourceOverrideConfig)
	if err != nil {
		err = errors2.Wrap(err, "failed to parse configuration override YAML")
		utils.HandleSpinnerError(spinner, sm, err)
		return err
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

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully applied one-off patch")

	// Search for the instance to get updated details
	searchRes, err := dataaccess.SearchInventory(cmd.Context(), token, fmt.Sprintf("resourceinstance:%s", instanceID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if len(searchRes.ResourceInstanceResults) == 0 {
		err = errors.New("failed to find the patched instance")
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