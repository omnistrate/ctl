package account

import (
	"fmt"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	createExample = `  # Create aws account
  omctl account create <name> --aws-account-id <aws-account-id>

  # Create gcp account
  omctl account create <name> --gcp-project-id <gcp-project-id> --gcp-project-number <gcp-project-number>`
)

var createCmd = &cobra.Command{
	Use:          "create [flags]",
	Short:        "Create an account",
	Long:         `Create an account with the specified name and cloud provider details.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	createCmd.Flags().String("aws-account-id", "", "AWS account ID")
	createCmd.Flags().String("gcp-project-id", "", "GCP project ID")
	createCmd.Flags().String("gcp-project-number", "", "GCP project number")

	err := createCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}

	createCmd.MarkFlagsMutuallyExclusive("aws-account-id", "gcp-project-id")
	createCmd.MarkFlagsOneRequired("aws-account-id", "gcp-project-id")
	createCmd.MarkFlagsRequiredTogether("gcp-project-id", "gcp-project-number")
}

func runCreate(cmd *cobra.Command, args []string) error {
	// Get flags
	awsAccountID, _ := cmd.Flags().GetString("aws-account-id")
	gcpProjectID, _ := cmd.Flags().GetString("gcp-project-id")
	gcpProjectNumber, _ := cmd.Flags().GetString("gcp-project-number")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Create account
	request := &accountconfigapi.CreateAccountConfigRequest{
		Token: token,
		Name:  args[0],
	}

	if awsAccountID != "" {
		// Get aws cloud provider id
		cloudProviderID, err := dataaccess.GetCloudProviderByName(token, "aws")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		request.CloudProviderID = accountconfigapi.CloudProviderID(cloudProviderID)
		request.AwsAccountID = &awsAccountID
		request.AwsBootstrapRoleARN = commonutils.ToPtr("arn:aws:iam::" + awsAccountID + ":role/omnistrate-bootstrap-role")
		request.Description = "AWS Account" + awsAccountID
	} else {
		// Get organization id
		user, err := dataaccess.DescribeUser(token)
		if err != nil {
			utils.PrintError(err)
			return err
		}

		// Get gcp cloud provider id
		cloudProviderID, err := dataaccess.GetCloudProviderByName(token, "gcp")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		request.CloudProviderID = accountconfigapi.CloudProviderID(cloudProviderID)
		request.GcpProjectID = &gcpProjectID
		request.GcpProjectNumber = &gcpProjectNumber
		request.GcpServiceAccountEmail = commonutils.ToPtr(fmt.Sprintf("bootstrap-%s@%s.iam.gserviceaccount.com", user.OrgID, gcpProjectID))
		request.Description = "GCP Account" + gcpProjectID
	}

	accountConfigId, err := dataaccess.CreateAccount(request)
	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess("Account created successfully")

	account, err := dataaccess.DescribeAccount(token, string(accountConfigId))
	if err != nil {
		utils.PrintError(err)
		return err
	}
	dataaccess.PrintNextStepVerifyAccountMsg(account)

	return nil
}
