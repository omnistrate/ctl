package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/config"
)

func DescribeResourceInstance(token string, serviceID, environmentID, instanceID string) (*inventoryapi.ResourceInstance, error) {
	instance, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
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

func RestartResourceInstance(token string, serviceID, environmentID, resourceID, instanceID string) error {
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

	err = instance.RestartResourceInstance(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}

func StartResourceInstance(token string, serviceID, environmentID, resourceID, instanceID string) error {
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

	err = instance.StartResourceInstance(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}

func StopResourceInstance(token string, serviceID, environmentID, resourceID, instanceID string) error {
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

	err = instance.StopResourceInstance(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}

func UpdateResourceInstance(token string, request inventoryapi.FleetUpdateResourceInstanceRequest) error {
	instance, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return err
	}

	request.Token = token

	err = instance.UpdateResourceInstance(context.Background(), &request)
	if err != nil {
		return err
	}

	return nil
}
