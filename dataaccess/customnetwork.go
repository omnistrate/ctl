package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/ctl/config"
)

func CreateCustomNetwork(token string, request customnetworkapi.CreateCustomNetworkRequest) (
	customNetwork *customnetworkapi.CustomNetwork, err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.CreateCustomNetwork(context.Background(), &request)
}

func DescribeCustomNetwork(token string, request customnetworkapi.DescribeCustomNetworkRequest) (
	customNetwork *customnetworkapi.CustomNetwork, err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.DescribeCustomNetwork(context.Background(), &request)
}

func ListCustomNetworks(token string, request customnetworkapi.ListCustomNetworksRequest) (
	customNetwork *customnetworkapi.ListCustomNetworksResult, err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.ListCustomNetworks(context.Background(), &request)
}

func DeleteCustomNetwork(token string, request customnetworkapi.DeleteCustomNetworkRequest) (
	err error) {
	var service *httpclientwrapper.CustomNetwork
	service, err = httpclientwrapper.NewCustomNetwork(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request.Token = token
	return service.DeleteCustomNetwork(context.Background(), &request)
}
