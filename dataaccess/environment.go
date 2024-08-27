package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	"strings"
)

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
)

func CreateServiceEnvironment(token string, request serviceenvironmentapi.CreateServiceEnvironmentRequest) (serviceenvironmentapi.ServiceEnvironmentID, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	request.Token = token

	res, err := service.CreateServiceEnvironment(context.Background(), &request)
	if err != nil {
		return "", err
	}
	return res, nil
}

func DescribeServiceEnvironment(token, serviceID, serviceEnvironmentID string) (*serviceenvironmentapi.DescribeServiceEnvironmentResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceenvironmentapi.DescribeServiceEnvironmentRequest{
		Token:     token,
		ServiceID: serviceenvironmentapi.ServiceID(serviceID),
		ID:        serviceenvironmentapi.ServiceEnvironmentID(serviceEnvironmentID),
	}

	res, err := service.DescribeServiceEnvironment(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func ListServiceEnvironments(token, serviceID string) (*serviceenvironmentapi.ListServiceEnvironmentsResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceenvironmentapi.ListServiceEnvironmentsRequest{
		Token:     token,
		ServiceID: serviceenvironmentapi.ServiceID(serviceID),
	}

	res, err := service.ListServiceEnvironment(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PromoteServiceEnvironment(token, serviceID, serviceEnvironmentID string) error {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := serviceenvironmentapi.PromoteServiceEnvironmentRequest{
		Token:     token,
		ServiceID: serviceenvironmentapi.ServiceID(serviceID),
		ID:        serviceenvironmentapi.ServiceEnvironmentID(serviceEnvironmentID),
	}

	err = service.PromoteServiceEnvironment(context.Background(), &request)
	if err != nil {
		return err
	}
	return nil
}

func PromoteServiceEnvironmentStatus(token, serviceID, serviceEnvironmentID string) (serviceenvironmentapi.PromoteServiceEnvironmentStatusResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := serviceenvironmentapi.PromoteServiceEnvironmentStatusRequest{
		Token:     token,
		ServiceID: serviceenvironmentapi.ServiceID(serviceID),
		ID:        serviceenvironmentapi.ServiceEnvironmentID(serviceEnvironmentID),
	}

	res, err := service.PromoteServiceEnvironmentStatus(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteServiceEnvironment(token, serviceID, serviceEnvironmentID string) error {
	service, err := httpclientwrapper.NewServiceEnvironment(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return err
	}

	request := serviceenvironmentapi.DeleteServiceEnvironmentRequest{
		Token:     token,
		ServiceID: serviceenvironmentapi.ServiceID(serviceID),
		ID:        serviceenvironmentapi.ServiceEnvironmentID(serviceEnvironmentID),
	}

	err = service.DeleteServiceEnvironment(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}

func FindEnvironment(token, serviceID, environmentType string) (*serviceenvironmentapi.DescribeServiceEnvironmentResult, error) {
	listRes, err := ListServiceEnvironments(token, serviceID)
	if err != nil {
		return nil, err
	}

	for _, id := range listRes.Ids {
		descRes, err := DescribeServiceEnvironment(token, serviceID, string(id))
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(string(descRes.Type), environmentType) {
			return descRes, nil
		}
	}

	return nil, ErrEnvironmentNotFound
}
