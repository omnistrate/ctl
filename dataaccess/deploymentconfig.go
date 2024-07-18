package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	deploymentconfigapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/deployment_config_api"
	"github.com/omnistrate/ctl/utils"
)

func GetDefaultDeploymentConfigID(token string) (deploymentconfigapi.DeploymentConfigID, error) {
	service, err := httpclientwrapper.NewDeploymentConfig(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	request := deploymentconfigapi.DescribeDeploymentConfigRequest{
		Token: token,
		ID:    "default",
	}

	res, err := service.DescribeDeploymentConfig(context.Background(), &request)
	if err != nil {
		return "", err
	}
	return res.ID, nil
}
