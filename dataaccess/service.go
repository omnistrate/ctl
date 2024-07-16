package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
)

func ListServices(token string) (*serviceapi.ListServiceResult, error) {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceapi.List{
		Token: token,
	}

	res, err := service.ListService(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DescribeService(serviceId, token string) (*serviceapi.DescribeServiceResult, error) {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceapi.DescribeServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceId),
	}

	res, err := service.DescribeService(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteService(serviceId, token string) error {
	service, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := serviceapi.DeleteServiceRequest{
		Token: token,
		ID:    serviceapi.ServiceID(serviceId),
	}

	err = service.DeleteService(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}
