package dataaccess

import (
	"context"
	"net/http"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

func DescribeHostCluster(ctx context.Context, token string, hostClusterID string) (*openapiclientfleet.HostCluster, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.HostclusterApiAPI.HostclusterApiDescribeHostCluster(ctxWithToken, hostClusterID)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	hostCluster, r, err := req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return hostCluster, nil
}

func ListHostClusters(ctx context.Context, token string, accountConfigID *string, regionID *string) (*openapiclientfleet.ListHostClustersResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.HostclusterApiAPI.HostclusterApiListHostClusters(ctxWithToken)

	if accountConfigID != nil {
		req = req.AccountConfigId(*accountConfigID)
	}
	if regionID != nil {
		req = req.RegionId(*regionID)
	}

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	hostClusters, r, err := req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return hostClusters, nil
}

func AdoptHostCluster(ctx context.Context, token string, hostClusterID string, cloudProvider string, region string, description string, userEmail *string) (*openapiclientfleet.AdoptHostClusterResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	adoptRequest := openapiclientfleet.AdoptHostClusterRequest2{
		CloudProvider: cloudProvider,
		Region:        region,
		Description:   description,
		Id:            hostClusterID,
	}

	if userEmail != nil && *userEmail != "" {
		adoptRequest.CustomerEmail = userEmail
	}

	req := apiClient.HostclusterApiAPI.HostclusterApiAdoptHostCluster(ctxWithToken).AdoptHostClusterRequest2(adoptRequest)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	hostCluster, r, err := req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return hostCluster, nil
}

func DeleteHostCluster(ctx context.Context, token string, hostClusterID string) error {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.HostclusterApiAPI.HostclusterApiDeleteHostCluster(ctxWithToken, hostClusterID)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err := req.Execute()
	if err != nil {
		return handleFleetError(err)
	}

	return nil
}
