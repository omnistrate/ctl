package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
)

func SearchInventory(token, query string) (*inventoryapi.SearchInventoryResult, error) {
	inventory, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.SearchInventory(context.Background(), &inventoryapi.SearchInventoryRequest{
		Token: token,
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListServiceOfferings(token, orgID string) (*inventoryapi.InventoryListServiceOfferingsResult, error) {
	inventory, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.ListServiceOffering(context.Background(), &inventoryapi.InventoryListServiceOfferingsRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DescribeServiceOfferingResource(token, serviceID, resourceID, instanceID, productTierID, productTierVersion string) (*inventoryapi.InventoryDescribeServiceOfferingResourceResult, error) {
	inventory, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.DescribeServiceOfferingResource(context.Background(), &inventoryapi.InventoryDescribeServiceOfferingResourceRequest{
		Token:              token,
		ServiceID:          inventoryapi.ServiceID(serviceID),
		ResourceID:         inventoryapi.ResourceID(resourceID),
		InstanceID:         instanceID,
		ProductTierID:      (*inventoryapi.ProductTierID)(utils.ToPtr(productTierID)),
		ProductTierVersion: utils.ToPtr(productTierVersion),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DescribeServiceOffering(token, serviceID, productTierID, productTierVersion string) (*inventoryapi.InventoryDescribeServiceOfferingResult, error) {
	inventory, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.DescribeServiceOffering(context.Background(), &inventoryapi.InventoryDescribeServiceOfferingRequest{
		Token:              token,
		ServiceID:          inventoryapi.ServiceID(serviceID),
		ProductTierID:      (*inventoryapi.ProductTierID)(utils.ToPtr(productTierID)),
		ProductTierVersion: utils.ToPtr(productTierVersion),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
