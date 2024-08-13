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
