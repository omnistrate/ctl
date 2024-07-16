package account

import (
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

var (
	accountExample = `
		# Create aws account
		create account <name> --aws-account-id <aws-account-id> --aws-bootstrap-role-arn <aws-bootstrap-role-arn>

		# Create gcp account
		omnistrate-ctl create account <name> --gcp-project-id <gcp-project-id> --gcp-project-number <gcp-project-number> --gcp-service-account-email <gcp-service-account-email>`
)

var (
	description            string
	awsAccountID           string
	awsBootstrapRoleARN    string
	gcpProjectID           string
	gcpProjectNumber       string
	gcpServiceAccountEmail string
)

// AccountCmd represents the create command
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

	AccountCmd.Flags().StringVarP(&description, "description", "", "", "Provide a description for the account")
	AccountCmd.Flags().StringVarP(&awsAccountID, "aws-account-id", "", "", "AWS account ID")
	AccountCmd.Flags().StringVarP(&awsBootstrapRoleARN, "aws-bootstrap-role-arn", "", "", "AWS bootstrap role ARN")
	AccountCmd.Flags().StringVarP(&gcpProjectID, "gcp-project-id", "", "", "GCP project ID")
	AccountCmd.Flags().StringVarP(&gcpProjectNumber, "gcp-project-number", "", "", "GCP project number")
	AccountCmd.Flags().StringVarP(&gcpServiceAccountEmail, "gcp-service-account-email", "", "", "GCP service account email")

	err := AccountCmd.MarkFlagRequired("name")
	if err != nil {
		return
	}

	AccountCmd.MarkFlagsMutuallyExclusive("aws-account-id", "gcp-project-id")
	AccountCmd.MarkFlagsOneRequired("aws-account-id", "gcp-project-id")
	AccountCmd.MarkFlagsRequiredTogether("aws-account-id", "aws-bootstrap-role-arn")
	AccountCmd.MarkFlagsRequiredTogether("gcp-project-id", "gcp-project-number", "gcp-service-account-email")
}

func run(cmd *cobra.Command, args []string) error {
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

	if description != "" {
		request.Description = description
	} else {
		request.Description = args[0]
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
		request.AwsBootstrapRoleARN = &awsBootstrapRoleARN
	} else {
		// Get gcp cloud provider id
		cloudProviderID, err := dataaccess.GetCloudProviderByName(token, "gcp")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		request.CloudProviderID = accountconfigapi.CloudProviderID(cloudProviderID)
		request.GcpProjectID = &gcpProjectID
		request.GcpProjectNumber = &gcpProjectNumber
		request.GcpServiceAccountEmail = &gcpServiceAccountEmail
	}

	_, err = dataaccess.CreateAccount(request)

	if err != nil {
		utils.PrintError(err)
		return err
	}
	utils.PrintSuccess("Account created successfully")

	return nil
}
