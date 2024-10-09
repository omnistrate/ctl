package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	producttierapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/product_tier_api"
	serviceapiapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_apiapi"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/utils"
)

func DeleteProductTier(ctx context.Context, token, serviceID, productTierID string) (err error) {
	service, err := httpclientwrapper.NewProductTier(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.DeleteProductTierRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceID),
		ID:        producttierapi.ProductTierID(productTierID),
	}

	if err = service.DeleteProductTier(context.Background(), request); err != nil {
		return
	}

	return
}

func DescribeProductTier(ctx context.Context, token, serviceID, productTierID string) (productTier *producttierapi.DescribeProductTierResult, err error) {
	service, err := httpclientwrapper.NewProductTier(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.DescribeProductTierRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceID),
		ID:        producttierapi.ProductTierID(productTierID),
	}

	productTier, err = service.DescribeProductTier(context.Background(), request)
	if err != nil {
		return
	}

	return
}

func ListProductTiers(ctx context.Context, token, serviceID string) (productTiers []*producttierapi.DescribeProductTierResult, err error) {
	service, err := httpclientwrapper.NewProductTier(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.ListProductTiersRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceID),
	}

	productTierIDs, err := service.ListProductTier(context.Background(), request)
	if err != nil {
		return
	}

	for _, productTierID := range productTierIDs.Ids {
		productTier, err := DescribeProductTier(ctx, token, serviceID, string(productTierID))
		if err != nil {
			return nil, err
		}
		productTiers = append(productTiers, productTier)
	}

	return
}

func ReleaseServicePlan(ctx context.Context, token, serviceID, serviceAPIID, productTierID string, versionSetName *string, isPreferred bool) (err error) {
	serviceApi, err := httpclientwrapper.NewServiceAPI(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &serviceapiapi.ReleaseServiceAPIRequest{
		Token:          token,
		ServiceID:      serviceapiapi.ServiceID(serviceID),
		ID:             serviceapiapi.ServiceAPIID(serviceAPIID),
		ProductTierID:  utils.ToPtr(serviceapiapi.ProductTierID(productTierID)),
		VersionSetName: versionSetName,
		VersionSetType: utils.ToPtr("Major"),
		IsPreferred:    isPreferred,
	}

	if err = serviceApi.ReleaseServiceAPI(context.Background(), request); err != nil {
		return
	}

	return
}

func DescribePendingChanges(ctx context.Context, token, serviceID, serviceAPIID, productTierID string) (pendingChanges *serviceapiapi.DescribePendingChangesResult, err error) {
	serviceApi, err := httpclientwrapper.NewServiceAPI(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &serviceapiapi.DescribePendingChangesRequest{
		Token:         token,
		ServiceID:     serviceapiapi.ServiceID(serviceID),
		ID:            serviceapiapi.ServiceAPIID(serviceAPIID),
		ProductTierID: utils.ToPtr(serviceapiapi.ProductTierID(productTierID)),
	}

	pendingChanges, err = serviceApi.DescribePendingChanges(context.Background(), request)
	if err != nil {
		return
	}

	return
}
