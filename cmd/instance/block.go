package instance

import (
	"errors"
	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	blockExample = `# Block an instance deployment
omctl instance block instance-abcd1234 --deployment-type terraform --deployment-name terraform-entity-name`

	TerraformDeploymentType DeploymentType = "terraform"
)

type DeploymentType string

var blockCmd = &cobra.Command{
	Use:          "block [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name>",
	Short:        "Block an instance deployment for your service",
	Long:         `This command helps you block the instance for your service.`,
	Example:      blockExample,
	RunE:         runBlock,
	SilenceUsage: true,
}

func init() {
	blockCmd.Flags().String("deployment-type", "", "Deployment type")
	blockCmd.Flags().String("deployment-name", "", "Deployment name")

	blockCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	blockCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	var err error
	if err = blockCmd.MarkFlagRequired("deployment-type"); err != nil {
		return
	}
	if err = blockCmd.MarkFlagRequired("deployment-name"); err != nil {
		return
	}
}

func runBlock(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	instanceID := args[0]

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
		msg := "Blocking instance..."
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

	// Block instance
	err = dataaccess.BlockResourceInstance(cmd.Context(), token, serviceID, environmentID, instanceID)
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

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully blocked instance deployment")
	// Print output
	err = utils.PrintTextTableJsonOutput(output, deploymentEntity)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
