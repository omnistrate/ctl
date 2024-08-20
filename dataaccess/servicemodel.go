package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	servicemodelapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_model_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeServiceModel(token, serviceId, serviceModelId string) (serviceModel *servicemodelapi.DescribeServiceModelResult, err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &servicemodelapi.DescribeServiceModelRequest{
		Token:     token,
		ServiceID: servicemodelapi.ServiceID(serviceId),
		ID:        servicemodelapi.ServiceModelID(serviceModelId),
	}

	serviceModel, err = fleetService.DescribeServiceModel(context.Background(), request)
	if err != nil {
		return
	}

	return
}
