package dataaccess

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	accountconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/account_config_api"
	"github.com/omnistrate/ctl/utils"
)

func ListAccounts(token string, cloudProvider string) (*accountconfigapi.ListAccountConfigResult, error) {
	account, err := httpclientwrapper.NewAccountConfig(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := accountconfigapi.ListAccountConfigRequest{
		Token:             token,
		CloudProviderName: accountconfigapi.CloudProvider(cloudProvider),
	}

	res, err := account.ListAccountConfig(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteAccount(accountConfigId, token string) error {
	service, err := httpclientwrapper.NewAccountConfig(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := accountconfigapi.DeleteAccountConfigRequest{
		Token: token,
		ID:    accountconfigapi.AccountConfigID(accountConfigId),
	}

	err = service.DeleteAccountConfig(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}

func CreateAccount(accountConfig *accountconfigapi.CreateAccountConfigRequest) (accountconfigapi.AccountConfigID, error) {
	service, err := httpclientwrapper.NewAccountConfig(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	res, err := service.CreateAccountConfig(context.Background(), accountConfig)
	if err != nil {
		return "", err
	}
	return res, nil
}

const warningMsgTemplate = `
WARNING! Account %s(%s) not verified. To complete the account configuration setup, follow the instructions below:
- For AWS CloudFormation users: Please create your CloudFormation Stack using the provided template at %s. Watch the CloudFormation guide at %s for help.
- For AWS/GCP Terraform users: Execute the Terraform scripts available at %s, by using the Account Config Identity ID below. For guidance our Terraform instructional video is at %s.`

func VerifyAccount() {
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return
	}

	// List all accounts
	listRes, err := ListAccounts(token, "all")
	if err != nil {
		utils.PrintError(err)
		return
	}

	// Warn if any accounts are not verified
	for _, account := range listRes.AccountConfigs {
		if account.Status != "READY" {
			awsCloudFormationTemplateURL := ""
			if account.AwsCloudFormationTemplateURL != nil {
				awsCloudFormationTemplateURL = *account.AwsCloudFormationTemplateURL
			}

			awsCloudFormationGuideURL := "https://youtu.be/Mu-4jppldwk"
			awsGcpTerraformScriptsURL := "https://github.com/omnistrate/account-setup"
			awsGcpTerraformGuideURL := "https://youtu.be/eKktc4QKgaA"

			var targetAccountID string
			if account.AwsAccountID != nil {
				targetAccountID = *account.AwsAccountID
			} else {
				targetAccountID = *account.GcpProjectID
			}

			utils.PrintWarning(fmt.Sprintf(warningMsgTemplate, account.Name, targetAccountID, awsCloudFormationTemplateURL,
				awsCloudFormationGuideURL, awsGcpTerraformScriptsURL, awsGcpTerraformGuideURL))
		}
	}
}
