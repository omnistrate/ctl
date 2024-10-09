package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/ctl/internal/config"
)

func CreateCustomNetwork(ctx context.Context, token string, request customnetworkapi.CreateCustomNetworkRequest) (
	customNetwork *customnetworkapi.CustomNetwork, err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.CreateCustomNetwork(ctx, &request)
}

func DescribeCustomNetwork(ctx context.Context, token string, request customnetworkapi.DescribeCustomNetworkRequest) (
	customNetwork *customnetworkapi.CustomNetwork, err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.DescribeCustomNetwork(ctx, &request)
}

func ListCustomNetworks(ctx context.Context, token string, request customnetworkapi.ListCustomNetworksRequest) (
	customNetwork *customnetworkapi.ListCustomNetworksResult, err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.ListCustomNetworks(ctx, &request)
}

func DeleteCustomNetwork(ctx context.Context, token string, request customnetworkapi.DeleteCustomNetworkRequest) (
	err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.DeleteCustomNetwork(ctx, &request)
}
