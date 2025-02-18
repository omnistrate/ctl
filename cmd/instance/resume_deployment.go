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
	resumeDeploymentExample = `# Resume an instance deployment
omctl instance resume-deployment instance-abcd1234 --deployment-type terraform --deployment-name my-terraform-deployment --deployment-action apply`
)

var resumeDeploymentCmd = &cobra.Command{
	Use:          "resume-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name> --deployment-action <deployment-action>",
	Short:        "Resume an instance deployment",
	Long:         `This command helps you resume the instance deployment.`,
	Example:      resumeDeploymentExample,
	RunE:         runResumeDeployment,
	SilenceUsage: true,
}

func init() {
	resumeDeploymentCmd.Flags().StringP("deployment-type", "t", "", "Deployment type")
	resumeDeploymentCmd.Flags().StringP("deployment-name", "n", "", "Deployment name")
	resumeDeploymentCmd.Flags().StringP("entity-action", "e", "", "Entity action")

	resumeDeploymentCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	resumeDeploymentCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	var err error
	if err = resumeDeploymentCmd.MarkFlagRequired("deployment-type"); err != nil {
		return
	}
	if err = resumeDeploymentCmd.MarkFlagRequired("deployment-name"); err != nil {
		return
	}
	if err = resumeDeploymentCmd.MarkFlagRequired("entity-action"); err != nil {
		return
	}
}

func runResumeDeployment(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

	// Retrieve args
	instanceID := args[0]
	// Retrieve flags
	deploymentType, err := cmd.Flags().GetString("deployment-type")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	deploymentName, err := cmd.Flags().GetString("deployment-name")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if deploymentName == "" {
		err = errors.New("deployment name is required")
		utils.PrintError(err)
		return err
	}

	var entityAction string
	if deploymentType == string(TerraformDeploymentType) {
		entityAction, err = cmd.Flags().GetString("entity-action")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		if entityAction == "" {
			err = errors.New("entity action is required")
			utils.PrintError(err)
			return err
		}
	}

	if deploymentType != string(TerraformDeploymentType) {
		err = errors.New("only terraform deployment type is supported")
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

	// Get instance deployment
	_, err = dataaccess.GetInstanceDeploymentEntity(cmd.Context(), token, instanceID, deploymentType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	// Resume instance deployment
	err = dataaccess.ResumeInstanceDeploymentEntity(cmd.Context(), token, instanceID, deploymentType, deploymentName, entityAction)
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

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully enabled override for instance deployment")
	// Print output
	err = utils.PrintTextTableJsonOutput(output, deploymentEntity)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil

}
