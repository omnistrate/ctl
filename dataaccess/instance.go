package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeInstance(token, serviceID, serviceEnvironmentID, instanceID string) (instance *inventoryapi.ResourceInstance, err error) {
	fleetService, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &inventoryapi.DescribeResourceInstanceRequestInternal{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(serviceEnvironmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
	}

	if instance, err = fleetService.DescribeResourceInstance(context.Background(), request); err != nil {
		return
	}

	return
}

func CreateInstance(token string, request inventoryapi.FleetCreateResourceInstanceRequest) (res *inventoryapi.CreateResourceInstanceResult, err error) {
	request.Token = token

	fleetService, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	if res, err = fleetService.CreateResourceInstance(context.Background(), &request); err != nil {
		return
	}

	return
}

func DeleteInstance(token, serviceID, serviceEnvironmentID, resourceID, instanceID string) (err error) {
	fleetService, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
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

	if err = fleetService.DeleteResourceInstance(context.Background(), request); err != nil {
		return
	}

	return
}
