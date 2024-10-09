package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/utils"
)

func SearchInventory(ctx context.Context, token, query string) (*inventoryapi.SearchInventoryResult, error) {
	inventory, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.SearchInventory(ctx, &inventoryapi.SearchInventoryRequest{
		Token: token,
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListServiceOfferings(ctx context.Context, token, orgID string) (*inventoryapi.InventoryListServiceOfferingsResult, error) {
	inventory, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.ListServiceOffering(ctx, &inventoryapi.InventoryListServiceOfferingsRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func DescribeServiceOfferingResource(ctx context.Context, token, serviceID, resourceID, instanceID, productTierID, productTierVersion string) (*inventoryapi.InventoryDescribeServiceOfferingResourceResult, error) {
	inventory, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.DescribeServiceOfferingResource(ctx, &inventoryapi.InventoryDescribeServiceOfferingResourceRequest{
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

func DescribeServiceOffering(ctx context.Context, token, serviceID, productTierID, productTierVersion string) (*inventoryapi.InventoryDescribeServiceOfferingResult, error) {
	inventory, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.DescribeServiceOffering(ctx, &inventoryapi.InventoryDescribeServiceOfferingRequest{
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
