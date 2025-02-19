package dataaccess

import (
	"context"
	"net/http"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

func ListServiceOfferings(ctx context.Context, token, orgID string) (inventory *openapiclientfleet.InventoryListServiceOfferingsResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiListServiceOfferings(ctxWithToken)
	req = req.OrgId(orgID)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	inventory, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return inventory, nil
}

func DescribeServiceOfferingResource(ctx context.Context, token, serviceID, resourceID, instanceID, productTierID, productTierVersion string) (res *openapiclientfleet.InventoryDescribeServiceOfferingResourceResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDescribeServiceOfferingResource(ctxWithToken, serviceID, resourceID, instanceID)
	req = req.ProductTierId(productTierID)
	req = req.ProductTierVersion(productTierVersion)
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

	return res, nil
}

func DescribeServiceOffering(ctx context.Context, token, serviceID, productTierID, productTierVersion string) (res *openapiclientfleet.InventoryDescribeServiceOfferingResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDescribeServiceOffering(ctxWithToken, serviceID)
	req = req.ProductTierId(productTierID)
	req = req.ProductTierVersion(productTierVersion)
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

	return res, nil
}
