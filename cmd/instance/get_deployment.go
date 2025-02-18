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
	getDeploymentExample = `  # Get the deployment entity metadata of the instance
	  omctl instance get-deployment instance-abcd1234 --deployment-type terraform --deployment-name my-terraform-deployment`
)

var getDeploymentCmd = &cobra.Command{
	Use:          "get-deployment [instance-id] --deployment-type <deployment-type> --deployment-name <deployment-name>",
	Short:        "Get the deployment entity metadata of the instance",
	Long:         `This command helps you get the deployment entity metadata of the instance.`,
	Example:      getDeploymentExample,
	RunE:         runGetDeployment,
	SilenceUsage: true,
}

func init() {
	getDeploymentCmd.Flags().StringP("deployment-type", "t", "", "Deployment type")
	getDeploymentCmd.Flags().StringP("deployment-name", "n", "", "Deployment name")

	getDeploymentCmd.Args = cobra.ExactArgs(1) // Require exactly one argument
	getDeploymentCmd.Flags().StringP("output", "o", "json", "Output format. Only json is supported")

	var err error
	if err = getDeploymentCmd.MarkFlagRequired("deployment-type"); err != nil {
		return
	}
	if err = getDeploymentCmd.MarkFlagRequired("deployment-name"); err != nil {
		return
	}
}

func runGetDeployment(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	if len(args) == 0 {
		err := errors.New("instance id is required")
		utils.PrintError(err)
		return err
	}

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
		msg := "Getting deployment entity metadata..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	deploymentEntity, err := dataaccess.GetInstanceDeploymentEntity(cmd.Context(), token, instanceID, deploymentType, deploymentName)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully got deployment entity metadata")
	// Print output
	err = utils.PrintTextTableJsonOutput(output, deploymentEntity)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
