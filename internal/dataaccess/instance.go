package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
	openapiclientfleet "github.com/omnistrate/omnistrate-sdk-go/fleet"
)

func DescribeInstance(ctx context.Context, token, serviceID, serviceEnvironmentID, instanceID string) (instance *openapiclientfleet.ResourceInstance, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)

	configuration := openapiclientfleet.NewConfiguration()
	apiClient := openapiclientfleet.NewAPIClient(configuration)
	res, r, err := apiClient.InventoryApiAPI.InventoryApiDescribeResourceInstance(
		ctxWithToken,
		serviceID,
		serviceEnvironmentID,
		instanceID,
	).Execute()

	err = handleFleetError(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return res, nil
}

func CreateInstance(ctx context.Context, token string, request inventoryapi.FleetCreateResourceInstanceRequest) (res *inventoryapi.CreateResourceInstanceResult, err error) {
	request.Token = token

	fleetService, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	if res, err = fleetService.CreateResourceInstance(ctx, &request); err != nil {
		return
	}

	return
}

func DeleteInstance(ctx context.Context, token, serviceID, serviceEnvironmentID, resourceID, instanceID string) (err error) {
	fleetService, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &inventoryapi.FleetDeleteResourceInstanceRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(serviceEnvironmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
		ResourceID:    inventoryapi.ResourceID(resourceID),
	}

	if err = fleetService.DeleteResourceInstance(ctx, request); err != nil {
		return
	}

	return
}
