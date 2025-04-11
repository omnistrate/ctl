package dataaccess

import (
	"context"
	"fmt"
	"strings"

	"github.com/omnistrate/ctl/internal/config"

	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/ctl/internal/utils"
)

func DescribeAccount(ctx context.Context, token string, id string) (*openapiclient.DescribeAccountConfigResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.AccountConfigApiAPI.AccountConfigApiDescribeAccountConfig(
		ctxWithToken,
		id,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return res, nil
}

func ListAccounts(ctx context.Context, token string, cloudProvider string) (*openapiclient.ListAccountConfigResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.AccountConfigApiAPI.AccountConfigApiListAccountConfig(
		ctxWithToken,
		cloudProvider,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return res, nil
}

func DeleteAccount(ctx context.Context, token, accountConfigID string) error {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	r, err := apiClient.AccountConfigApiAPI.AccountConfigApiDeleteAccountConfig(
		ctxWithToken,
		accountConfigID,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return err
	}

	r.Body.Close()
	return nil
}

func CreateAccount(ctx context.Context, token string, accountConfig openapiclient.CreateAccountConfigRequest2) (string, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.AccountConfigApiAPI.AccountConfigApiCreateAccountConfig(
		ctxWithToken,
	).CreateAccountConfigRequest2(accountConfig).Execute()

	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	r.Body.Close()
	return strings.Trim(res, "\"\n"), nil
}

const (
	AccountNotVerifiedWarningMsgTemplate = `
WARNING! Account %s(%s) is not verified. To complete the account configuration setup, follow the instructions below:
- For AWS CloudFormation users: Please create your CloudFormation Stack using the provided template at %s. Watch the CloudFormation guide at %s for help.
- For AWS/GCP Terraform users: Execute the Terraform scripts available at %s, by using the Account Config Identity ID below. For guidance our Terraform instructional video is at %s.
- For Azure users: Execute the Azure bootstrap script using the command provided in your account configuration details. This will set up the necessary Azure AD applications and role assignments.`

	NextStepVerifyAccountMsgTemplate = `
Next step:
Verify your account.

- For AWS CloudFormation users: Please create your CloudFormation Stack using the provided template at %s. Watch the CloudFormation guide at %s for help.
- For AWS/GCP Terraform users: Execute the Terraform scripts available at %s, by using the Account Config Identity ID below. For guidance our Terraform instructional video is at %s.`

	AwsCloudFormationGuideURL = "https://youtu.be/Mu-4jppldwk"
	AwsGcpTerraformScriptsURL = "https://github.com/omnistrate-oss/account-setup"
	AwsGcpTerraformGuideURL   = "https://youtu.be/eKktc4QKgaA"
)

func PrintNextStepVerifyAccountMsg(account *openapiclient.DescribeAccountConfigResult) {
	awsCloudFormationTemplateURL := ""
	if account.AwsCloudFormationTemplateURL != nil {
		awsCloudFormationTemplateURL = *account.AwsCloudFormationTemplateURL
	}

	fmt.Println(fmt.Sprintf(NextStepVerifyAccountMsgTemplate, awsCloudFormationTemplateURL,
		AwsCloudFormationGuideURL, AwsGcpTerraformScriptsURL, AwsGcpTerraformGuideURL))
}

func PrintAccountNotVerifiedWarning(account *openapiclient.DescribeAccountConfigResult) {
	awsCloudFormationTemplateURL := ""
	if account.AwsCloudFormationTemplateURL != nil {
		awsCloudFormationTemplateURL = *account.AwsCloudFormationTemplateURL
	}

	var targetAccountID string
	if account.AwsAccountID != nil {
		targetAccountID = *account.AwsAccountID
	} else if account.GcpProjectID != nil {
		targetAccountID = *account.GcpProjectID
	} else if account.AzureSubscriptionID != nil {
		targetAccountID = *account.AzureSubscriptionID
	} else {
		targetAccountID = "unknown"
	}

	name := account.Name
	if name == "" {
		name = "Unnamed Account"
	}

	utils.PrintWarning(fmt.Sprintf(AccountNotVerifiedWarningMsgTemplate, name, targetAccountID, awsCloudFormationTemplateURL,
		AwsCloudFormationGuideURL, AwsGcpTerraformScriptsURL, AwsGcpTerraformGuideURL))
}

func AskVerifyAccountIfAny(ctx context.Context) {
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return
	}

	// List all accounts
	listRes, err := ListAccounts(ctx, token, "all")
	if err != nil {
		utils.PrintError(err)
		return
	}

	// Warn if any accounts are not verified
	for _, account := range listRes.AccountConfigs {
		if account.Status != "READY" {
			PrintAccountNotVerifiedWarning(&account)
		}
	}
}
