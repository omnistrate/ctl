package dataaccess

import (
	"context"
	"net/http"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
)

func CreateResourceInstance(ctx context.Context, token string, request inventoryapi.FleetCreateResourceInstanceRequest) (res *inventoryapi.CreateResourceInstanceResult, err error) {
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

func DeleteResourceInstance(ctx context.Context, token, serviceID, serviceEnvironmentID, resourceID, instanceID string) (err error) {
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

func RestartResourceInstance(ctx context.Context, token string, serviceID, environmentID, resourceID, instanceID string) error {
	instance, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	request := inventoryapi.FleetRestartResourceInstanceRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
		ResourceID:    inventoryapi.ResourceID(resourceID),
	}

	err = instance.RestartResourceInstance(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}

func StartResourceInstance(ctx context.Context, token string, serviceID, environmentID, resourceID, instanceID string) error {
	instance, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	request := inventoryapi.FleetStartResourceInstanceRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
		ResourceID:    inventoryapi.ResourceID(resourceID),
	}

	err = instance.StartResourceInstance(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}

func StopResourceInstance(ctx context.Context, token string, serviceID, environmentID, resourceID, instanceID string) error {
	instance, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	request := inventoryapi.FleetStopResourceInstanceRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
		ResourceID:    inventoryapi.ResourceID(resourceID),
	}

	err = instance.StopResourceInstance(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}

func UpdateResourceInstance(ctx context.Context, token string, request inventoryapi.FleetUpdateResourceInstanceRequest) error {
	instance, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	request.Token = token

	err = instance.UpdateResourceInstance(ctx, &request)
	if err != nil {
		return err
	}

	return nil
}
