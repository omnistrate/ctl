package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	cloudproviderapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/cloud_provider_api"
	"github.com/omnistrate/ctl/config"
)

func GetCloudProviderByName(token string, cloudProvider string) (cloudproviderapi.CloudProviderID, error) {
	service, err := httpclientwrapper.NewCloudProvider(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return "", err
	}

	request := cloudproviderapi.GetCloudProviderByNameRequest{
		Token: token,
		Name:  cloudProvider,
	}

	res, err := service.GetCloudProviderByName(context.Background(), &request)
	if err != nil {
		return "", err
	}
	return res, nil
}
