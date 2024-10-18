package dataaccess

import (
	"context"
	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	"net/http"
)

func FleetDescribeCustomNetwork(
	ctx context.Context, token string, id string) (
	customNetwork *openapiclientfleet.FleetCustomNetwork, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.FleetCustomNetworkApiAPI.FleetCustomNetworkApiDescribeCustomNetwork(ctxWithToken, id)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	customNetwork, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return
}

func FleetListCustomNetworks(
	ctx context.Context, token string, cloudProviderName *string, cloudProviderRegion *string) (
	customNetworks *openapiclientfleet.FleetListCustomNetworksResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.FleetCustomNetworkApiAPI.FleetCustomNetworkApiListCustomNetworks(ctxWithToken)
	if cloudProviderName != nil {
		req.CloudProviderName(*cloudProviderName)
	}
	if cloudProviderRegion != nil {
		req.CloudProviderRegion(*cloudProviderRegion)
	}

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	customNetworks, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return
}
