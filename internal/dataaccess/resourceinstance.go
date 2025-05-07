package dataaccess

import (
	"context"
	"github.com/omnistrate/ctl/internal/utils"
	"net/http"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

func CreateResourceInstance(ctx context.Context, token string,
	serviceProviderId string, serviceKey string, serviceAPIVersion string, serviceEnvironmentKey string, serviceModelKey string, productTierKey string, resourceKey string,
	request openapiclientfleet.FleetCreateResourceInstanceRequest2) (res *openapiclientfleet.FleetCreateResourceInstanceResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiCreateResourceInstance(
		ctxWithToken,
		serviceProviderId,
		serviceKey,
		serviceAPIVersion,
		serviceEnvironmentKey,
		serviceModelKey,
		productTierKey,
		resourceKey,
	).FleetCreateResourceInstanceRequest2(request)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	res, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func RestoreResourceInstanceSnapshot(ctx context.Context, token string, serviceID, environmentID, snapshotID string, formattedParams map[string]any) (res *openapiclientfleet.FleetRestoreResourceInstanceResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiRestoreResourceInstanceFromSnapshot(
		ctxWithToken,
		serviceID,
		environmentID,
		snapshotID,
	).FleetRestoreResourceInstanceFromSnapshotRequest2(
		openapiclientfleet.FleetRestoreResourceInstanceFromSnapshotRequest2{
			InputParametersOverride: formattedParams,
			NetworkType:             utils.ToPtr("PUBLIC"), // We always set this to PUBLIC for now
		})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	res, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func DescribeResourceInstanceSnapshot(ctx context.Context, token string, serviceID, environmentID, instanceID, snapshotID string) (res *openapiclientfleet.FleetDescribeInstanceSnapshotResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDescribeResourceInstanceSnapshot(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
		snapshotID,
	)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	res, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func ListResourceInstanceSnapshots(ctx context.Context, token string, serviceID, environmentID, instanceID string) (res *openapiclientfleet.FleetListInstanceSnapshotResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiListResourceInstanceSnapshots(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	res, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func TriggerResourceInstanceAutoBackup(ctx context.Context, token string, serviceID, environmentID, instanceID string) (res *openapiclientfleet.FleetAutomaticInstanceSnapshotCreationResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiTriggerAutomaticResourceInstanceSnapshotCreation(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	res, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func DeleteResourceInstance(ctx context.Context, token, serviceID, environmentID, resourceID, instanceID string) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDeleteResourceInstance(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	).FleetDeleteResourceInstanceRequest2(openapiclientfleet.FleetDeleteResourceInstanceRequest2{
		ResourceId: resourceID,
	})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err = req.Execute()
	if err != nil {
		return handleFleetError(err)
	}
	return
}

func DescribeResourceInstance(ctx context.Context, token string, serviceID, environmentID, instanceID string) (resp *openapiclientfleet.ResourceInstance, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDescribeResourceInstance(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	resp, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func UpdateResourceInstanceDebugMode(ctx context.Context, token string, serviceID, environmentID, instanceID string, enable bool) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiUpdateResourceInstanceDebugMode(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	).FleetUpdateResourceInstanceDebugModeRequest2(openapiclientfleet.FleetUpdateResourceInstanceDebugModeRequest2{
		Enable: enable,
	})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err = req.Execute()
	if err != nil {
		return handleFleetError(err)
	}
	return
}

func RestartResourceInstance(ctx context.Context, token string, serviceID, environmentID, resourceID, instanceID string) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiRestartResourceInstance(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	).FleetRestartResourceInstanceRequest2(openapiclientfleet.FleetRestartResourceInstanceRequest2{
		ResourceId: resourceID,
	})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err = req.Execute()
	if err != nil {
		return handleFleetError(err)
	}
	return
}

func StartResourceInstance(ctx context.Context, token string, serviceID, environmentID, resourceID, instanceID string) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiStartResourceInstance(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	).FleetStartResourceInstanceRequest2(openapiclientfleet.FleetStartResourceInstanceRequest2{
		ResourceId: resourceID,
	})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err = req.Execute()
	if err != nil {
		return handleFleetError(err)
	}
	return
}

func StopResourceInstance(ctx context.Context, token string, serviceID, environmentID, resourceID, instanceID string) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiStopResourceInstance(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	).FleetStopResourceInstanceRequest2(openapiclientfleet.FleetStopResourceInstanceRequest2{
		ResourceId: resourceID,
	})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err = req.Execute()
	if err != nil {
		return handleFleetError(err)
	}
	return
}

func UpdateResourceInstance(ctx context.Context, token string,
	serviceID, environmentID, instanceID string,
	resourceId string, requestParameters map[string]any) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiUpdateResourceInstance(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	).FleetUpdateResourceInstanceRequest2(openapiclientfleet.FleetUpdateResourceInstanceRequest2{
		ResourceId:    resourceId,
		RequestParams: requestParameters,
	})

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	r, err = req.Execute()
	if err != nil {
		return handleFleetError(err)
	}
	return
}
