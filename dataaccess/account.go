package dataaccess

import (
	"context"
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
