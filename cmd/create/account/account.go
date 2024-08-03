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
	accountExample = `  # Create aws account
  create account <name> --aws-account-id <aws-account-id>

  # Create gcp account
  omnistrate-ctl create account <name> --gcp-project-id <gcp-project-id> --gcp-project-number <gcp-project-number>`
)

var AccountCmd = &cobra.Command{
	Use:          "account <name> [flags]",
	Short:        "Create a account",
	Long:         ``,
	Example:      accountExample,
	RunE:         run,
	SilenceUsage: true,
}

func init() {
	AccountCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	AccountCmd.Flags().String("aws-account-id", "", "AWS account ID")
	AccountCmd.Flags().String("gcp-project-id", "", "GCP project ID")
	AccountCmd.Flags().String("gcp-project-number", "", "GCP project number")

	err := AccountCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}

	AccountCmd.MarkFlagsMutuallyExclusive("aws-account-id", "gcp-project-id")
	AccountCmd.MarkFlagsOneRequired("aws-account-id", "gcp-project-id")
	AccountCmd.MarkFlagsRequiredTogether("gcp-project-id", "gcp-project-number")
}

func run(cmd *cobra.Command, args []string) error {
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
