package instance

import (
	"encoding/json"
	"errors"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	disableDebugModeExample = `# Disable instance deployment debug mode
omctl instance disable-debug-mode instance-abcd1234 --resource-name my-terraform-deployment --deployment-action apply`
)

var disableDebugModeCmd = &cobra.Command{
	Use:          "disable-debug-mode [instance-id] --resource-name <resource-name> --deployment-action <deployment-action>",
	Short:        "Disable instance debug mode",
	Long:         `This command helps you disable instance debug mode.`,
	Example:      disableDebugModeExample,
	RunE:         runDisableDebug,
	SilenceUsage: true,
}

func init() {
	disableDebugModeCmd.Flags().StringP("resource-name", "r", "", "Resource name")
	disableDebugModeCmd.Flags().StringP("deployment-action", "e", "", "Deployment action")

	disableDebugModeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	disableDebugModeCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	var err error
	if err = disableDebugModeCmd.MarkFlagRequired("resource-name"); err != nil {
		return
	}
	if err = disableDebugModeCmd.MarkFlagRequired("deployment-action"); err != nil {
		return
	}
}

func runDisableDebug(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]
	// Retrieve flags
	resourceName, err := cmd.Flags().GetString("resource-name")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if resourceName == "" {
		err = errors.New("resource name is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve flags
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
		msg := "Resuming deployment..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	resourceID, resourceType, err := getResourceFromInstance(cmd.Context(), token, instanceID, resourceName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	var deploymentAction string
	if resourceType == string(TerraformDeploymentType) {
		deploymentAction, err = cmd.Flags().GetString("deployment-action")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if deploymentAction == "" {
			err = errors.New("deployment action is required")
			utils.PrintError(err)
			return err
		}
	}

	if resourceType != string(TerraformDeploymentType) {
		err = errors.New("only terraform deployment type is supported")
		utils.PrintError(err)
		return err
	}

	var deploymentName string
	switch strings.ToLower(resourceType) {
	case string(TerraformDeploymentType):
		deploymentName = getTerraformDeploymentName(resourceID, instanceID)
	}

	// Get instance deployment
	_, err = dataaccess.GetInstanceDeploymentEntity(cmd.Context(), token, instanceID, resourceType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Resume instance deployment
	err = dataaccess.ResumeInstanceDeploymentEntity(cmd.Context(), token, instanceID, resourceType, deploymentName, deploymentAction)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.PrintWarning("The instance is currently locked for operations. Debug mode has been disabled for the deployment, but to fully unlock the instance and resume normal operations, you'll need to perform an instance upgrade.")

	// Describe deployment entity
	deploymentEntity, err := dataaccess.GetInstanceDeploymentEntity(cmd.Context(), token, instanceID, resourceType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	switch resourceType {
	case string(TerraformDeploymentType):
		// Parse JSON
		var response TerraformResponse
		err = json.Unmarshal([]byte(deploymentEntity), &response)
		if err != nil {
			utils.PrintError(errors2.Errorf("Error parsing instance deployment entity response: %v\n", err))
			return err
		}

		displayResource := TerraformResponse{}
		displayResource.Files = response.Files
		displayResource.Files.FilesContents = nil
		displayResource.SyncState = response.SyncState
		displayResource.SyncError = response.SyncError

		// Convert to JSON
		var displayOutput []byte
		displayOutput, err = json.MarshalIndent(displayResource, "", "  ")
		if err != nil {
			utils.PrintError(errors2.Errorf("Error converting instance deployment entity response to JSON: %v\n", err))
			return err
		}

		deploymentEntity = string(displayOutput)
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully enabled override for instance deployment")
	// Print output
	err = utils.PrintTextTableJsonOutput(output, deploymentEntity)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil

}
