package dataaccess

import (
	"context"
	"strings"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceenvironmentapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_environment_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/pkg/errors"
)

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
)

func CreateServiceEnvironment(ctx context.Context, token string, request serviceenvironmentapi.CreateServiceEnvironmentRequest) (serviceenvironmentapi.ServiceEnvironmentID, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(config.GetHostScheme(), config.GetHost())
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

func DescribeServiceEnvironment(ctx context.Context, token, serviceID, serviceEnvironmentID string) (*serviceenvironmentapi.DescribeServiceEnvironmentResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(config.GetHostScheme(), config.GetHost())
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

func ListServiceEnvironments(ctx context.Context, token, serviceID string) (*serviceenvironmentapi.ListServiceEnvironmentsResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(config.GetHostScheme(), config.GetHost())
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

func PromoteServiceEnvironment(ctx context.Context, token, serviceID, serviceEnvironmentID string) error {
	service, err := httpclientwrapper.NewServiceEnvironment(config.GetHostScheme(), config.GetHost())
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

func PromoteServiceEnvironmentStatus(ctx context.Context, token, serviceID, serviceEnvironmentID string) (serviceenvironmentapi.PromoteServiceEnvironmentStatusResult, error) {
	service, err := httpclientwrapper.NewServiceEnvironment(config.GetHostScheme(), config.GetHost())
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

func DeleteServiceEnvironment(ctx context.Context, token, serviceID, serviceEnvironmentID string) error {
	service, err := httpclientwrapper.NewServiceEnvironment(config.GetHostScheme(), config.GetHost())
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

func FindEnvironment(ctx context.Context, token, serviceID, environmentType string) (*serviceenvironmentapi.DescribeServiceEnvironmentResult, error) {
	listRes, err := ListServiceEnvironments(ctx, token, serviceID)
	if err != nil {
		return nil, err
	}

	for _, id := range listRes.Ids {
		descRes, err := DescribeServiceEnvironment(ctx, token, serviceID, string(id))
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(string(descRes.Type), environmentType) {
			return descRes, nil
		}
	}

	return nil, ErrEnvironmentNotFound
}
