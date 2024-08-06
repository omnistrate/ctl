package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeResourceInstance(token string, serviceID, environmentID, instanceID string) (*inventoryapi.ResourceInstance, error) {
	instance, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := inventoryapi.DescribeResourceInstanceRequestInternal{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		InstanceID:    inventoryapi.ResourceInstanceID(instanceID),
	}

	res, err := instance.DescribeResourceInstance(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
