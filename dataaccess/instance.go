package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeInstance(token, serviceId, serviceEnvironmentId, instanceId string) (instance *inventoryapi.ResourceInstance, err error) {
	fleetService, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &inventoryapi.DescribeResourceInstanceRequestInternal{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceId),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(serviceEnvironmentId),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceId),
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

func DeleteInstance(token, serviceId, serviceEnvironmentId, resourceId, instanceId string) (err error) {
	fleetService, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &inventoryapi.FleetDeleteResourceInstanceRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceId),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(serviceEnvironmentId),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceId),
		ResourceID:    inventoryapi.ResourceID(resourceId),
	}

	if err = fleetService.DeleteResourceInstance(context.Background(), request); err != nil {
		return
	}

	return
}
