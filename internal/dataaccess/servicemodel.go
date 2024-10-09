package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	servicemodelapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_model_api"
	"github.com/omnistrate/ctl/internal/config"
)

func DescribeServiceModel(ctx context.Context, token, serviceID, serviceModelID string) (serviceModel *servicemodelapi.DescribeServiceModelResult, err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &servicemodelapi.DescribeServiceModelRequest{
		Token:     token,
		ServiceID: servicemodelapi.ServiceID(serviceID),
		ID:        servicemodelapi.ServiceModelID(serviceModelID),
	}

	serviceModel, err = fleetService.DescribeServiceModel(ctx, request)
	if err != nil {
		return
	}

	return
}

func EnableServiceModelFeature(ctx context.Context, token, serviceID, serviceModelID, featureName string, featureConfiguration map[string]any) (err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &servicemodelapi.EnableServiceModelFeatureRequest{
		Token:         token,
		ServiceID:     servicemodelapi.ServiceID(serviceID),
		ID:            servicemodelapi.ServiceModelID(serviceModelID),
		Feature:       servicemodelapi.ServiceModelFeatureName(featureName),
		Configuration: featureConfiguration,
	}

	err = fleetService.EnableServiceModelFeature(ctx, request)
	if err != nil {
		return
	}

	return
}

func DisableServiceModelFeature(ctx context.Context, token, serviceID, serviceModelID, featureName string) (err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &servicemodelapi.DisableServiceModelFeatureRequest{
		Token:     token,
		ServiceID: servicemodelapi.ServiceID(serviceID),
		ID:        servicemodelapi.ServiceModelID(serviceModelID),
		Feature:   servicemodelapi.ServiceModelFeatureName(featureName),
	}

	err = fleetService.DisableServiceModelFeature(ctx, request)
	if err != nil {
		return
	}

	return
}
