package account

import (
	"fmt"

	"github.com/chelnak/ysmrr"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	createExample = `# Create aws account
omctl account create [account-name] --aws-account-id=[account-id]

# Create gcp account
omctl account create [account-name] --gcp-project-id=[project-id] --gcp-project-number=[project-number]`
)

var createCmd = &cobra.Command{
	Use:          "create [account-name] [--aws-account-id=account-id] [--gcp-project-id=project-id] [--gcp-project-number=project-number]",
	Short:        "Create a Cloud Provider Account",
	Long:         `This command helps you create a Cloud Provider Account in your account list.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	createCmd.Flags().String("aws-account-id", "", "AWS account ID")
	createCmd.Flags().String("gcp-project-id", "", "GCP project ID")
	createCmd.Flags().String("gcp-project-number", "", "GCP project number")

	// TODO: Uncomment the following lines to add validation to the flags
	// createCmd.MarkFlagsMutuallyExclusive("aws-account-id", "gcp-project-id")
	// createCmd.MarkFlagsOneRequired("aws-account-id", "gcp-project-id")
	// createCmd.MarkFlagsRequiredTogether("gcp-project-id", "gcp-project-number")
}

func runCreate(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve args
	var name string
	if len(args) > 0 {
		name = args[0]
	}

	// Retrieve flags
	awsAccountID, _ := cmd.Flags().GetString("aws-account-id")
	gcpProjectID, _ := cmd.Flags().GetString("gcp-project-id")
	gcpProjectNumber, _ := cmd.Flags().GetString("gcp-project-number")
	output, _ := cmd.Flags().GetString("output")

	// Validate user login
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Creating account..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Prepare request
	request := &accountconfigapi.CreateAccountConfigRequest{
		Token: token,
		Name:  name,
	}

	if awsAccountID != "" {
		// Get aws cloud provider id
		cloudProviderID, err := dataaccess.GetCloudProviderByName(cmd.Context(), token, "aws")
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		request.CloudProviderID = accountconfigapi.CloudProviderID(cloudProviderID)
		request.AwsAccountID = &awsAccountID
		request.AwsBootstrapRoleARN = utils.ToPtr("arn:aws:iam::" + awsAccountID + ":role/omnistrate-bootstrap-role")
		request.Description = "AWS Account" + awsAccountID
	} else {
		// Get organization id
		user, err := dataaccess.DescribeUser(cmd.Context(), token)
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		// Get gcp cloud provider id
		cloudProviderID, err := dataaccess.GetCloudProviderByName(cmd.Context(), token, "gcp")
		if err != nil {
			utils.HandleSpinnerError(spinner, sm, err)
			return err
		}

		request.CloudProviderID = accountconfigapi.CloudProviderID(cloudProviderID)
		request.GcpProjectID = &gcpProjectID
		request.GcpProjectNumber = &gcpProjectNumber
		request.GcpServiceAccountEmail = utils.ToPtr(fmt.Sprintf("bootstrap-%s@%s.iam.gserviceaccount.com", user.OrgId, gcpProjectID))
		request.Description = "GCP Account" + gcpProjectID
	}

	// Create account
	accountConfigID, err := dataaccess.CreateAccount(cmd.Context(), request)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}
	utils.HandleSpinnerSuccess(spinner, sm, "Successfully created account")

	// Describe account
	account, err := dataaccess.DescribeAccount(cmd.Context(), token, string(accountConfigID))
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	err = utils.PrintTextTableJsonOutput(output, account)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print next step
	if output != "json" {
		dataaccess.PrintNextStepVerifyAccountMsg(account)
	}

	return nil
}
