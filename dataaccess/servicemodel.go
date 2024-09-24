package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	servicemodelapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_model_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeServiceModel(token, serviceID, serviceModelID string) (serviceModel *servicemodelapi.DescribeServiceModelResult, err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &servicemodelapi.DescribeServiceModelRequest{
		Token:     token,
		ServiceID: servicemodelapi.ServiceID(serviceID),
		ID:        servicemodelapi.ServiceModelID(serviceModelID),
	}

	serviceModel, err = fleetService.DescribeServiceModel(context.Background(), request)
	if err != nil {
		return
	}

	return
}

func EnableServiceModelFeature(token, serviceID, serviceModelID, featureName string, featureConfiguration map[string]any) (err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(utils.GetHostScheme(), utils.GetHost())
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

	err = fleetService.EnableServiceModelFeature(context.Background(), request)
	if err != nil {
		return
	}

	return
}

func DisableServiceModelFeature(token, serviceID, serviceModelID, featureName string) (err error) {
	fleetService, err := httpclientwrapper.NewServiceModel(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &servicemodelapi.DisableServiceModelFeatureRequest{
		Token:     token,
		ServiceID: servicemodelapi.ServiceID(serviceID),
		ID:        servicemodelapi.ServiceModelID(serviceModelID),
		Feature:   servicemodelapi.ServiceModelFeatureName(featureName),
	}

	err = fleetService.DisableServiceModelFeature(context.Background(), request)
	if err != nil {
		return
	}

	return
}
