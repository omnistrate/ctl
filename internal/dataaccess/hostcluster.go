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
