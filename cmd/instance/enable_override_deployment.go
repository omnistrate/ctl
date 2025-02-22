package instance

import (
	"encoding/json"
	"errors"
	"github.com/chelnak/ysmrr"
	"github.com/cqroot/prompt"
	"github.com/cqroot/prompt/choose"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	enableOverrideExample = `# Enable override for an instance deployment
omctl instance enable-override-deployment <instance-id> --deployment-type terraform --deployment-name terraform-entity-name --force`
)

var enableOverrideCmd = &cobra.Command{
	Use:          "enable-override-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> --force",
	Short:        "Enable override for an instance deployment",
	Long:         `This command helps you enable override for an instance deployment`,
	Example:      enableOverrideExample,
	RunE:         runEnableOverride,
	SilenceUsage: true,
}

func init() {
	enableOverrideCmd.Flags().StringP("deployment-type", "t", "", "Deployment type")
	enableOverrideCmd.Flags().StringP("deployment-name", "n", "", "Deployment name")
	enableOverrideCmd.Flags().BoolP("force", "f", false, "Force enable override")

	enableOverrideCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	enableOverrideCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	var err error
	if err = enableOverrideCmd.MarkFlagRequired("deployment-type"); err != nil {
		return
	}
	if err = enableOverrideCmd.MarkFlagRequired("deployment-name"); err != nil {
		return
	}
}

func runEnableOverride(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]

	isForce, err := cmd.Flags().GetBool("force")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if !isForce {
		// Prompt user to confirm
		var choice string
		choice, err = prompt.New().Ask("Enable override will interrupt ongoing terraform operations, continue to proceed?").
			Choose([]string{
				"Yes",
				"No",
			}, choose.WithTheme(choose.ThemeArrow))
		if err != nil {
			utils.PrintError(err)
			return err
		}

		switch choice {
		case "Yes":
			break
		case "No":
			utils.PrintInfo("Operation cancelled")
			return nil
		}
	}

	// Retrieve flags
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

	// Retrieve flags
	deploymentType, err := cmd.Flags().GetString("deployment-type")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate deployment type
	if deploymentType != string(TerraformDeploymentType) {
		err = errors.New("only terraform deployment type is supported")
		utils.PrintError(err)
		return err
	}

	// Retrieve flags
	deploymentName, err := cmd.Flags().GetString("deployment-name")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate deployment name
	if deploymentName == "" {
		err = errors.New("deployment name is required")
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
		msg := "Enabling override for instance deployment..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	_, err = dataaccess.GetInstanceDeploymentEntity(cmd.Context(), token, instanceID, deploymentType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Check if instance exists
	serviceID, environmentID, _, _, err := getInstance(cmd.Context(), token, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Enable override
	err = dataaccess.EnableResourceInstanceOverride(cmd.Context(), token, serviceID, environmentID, instanceID)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Pause deployment
	err = dataaccess.PauseInstanceDeploymentEntity(cmd.Context(), token, instanceID, deploymentType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Describe deployment entity
	deploymentEntity, err := dataaccess.GetInstanceDeploymentEntity(cmd.Context(), token, instanceID, deploymentType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	switch deploymentType {
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
