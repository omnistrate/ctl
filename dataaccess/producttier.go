package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	producttierapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/product_tier_api"
	serviceapiapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_apiapi"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/utils"
)

func DeleteProductTier(token, serviceId, productTierId string) (err error) {
	service, err := httpclientwrapper.NewProductTier(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.DeleteProductTierRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceId),
		ID:        producttierapi.ProductTierID(productTierId),
	}

	if err = service.DeleteProductTier(context.Background(), request); err != nil {
		return
	}

	return
}

func DescribeProductTier(token, serviceId, productTierId string) (productTier *producttierapi.DescribeProductTierResult, err error) {
	service, err := httpclientwrapper.NewProductTier(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.DescribeProductTierRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceId),
		ID:        producttierapi.ProductTierID(productTierId),
	}

	productTier, err = service.DescribeProductTier(context.Background(), request)
	if err != nil {
		return
	}

	return
}

func ListProductTiers(token, serviceId string) (productTiers []*producttierapi.DescribeProductTierResult, err error) {
	service, err := httpclientwrapper.NewProductTier(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.ListProductTiersRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceId),
	}

	productTierIds, err := service.ListProductTier(context.Background(), request)
	if err != nil {
		return
	}

	for _, productTierId := range productTierIds.Ids {
		productTier, err := DescribeProductTier(token, serviceId, string(productTierId))
		if err != nil {
			return nil, err
		}
		productTiers = append(productTiers, productTier)
	}

	return
}

func ReleaseServicePlan(token, serviceId, serviceApiId, productTierId string, versionSetName *string, isPreferred bool) (err error) {
	serviceApi, err := httpclientwrapper.NewServiceAPI(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &serviceapiapi.ReleaseServiceAPIRequest{
		Token:          token,
		ServiceID:      serviceapiapi.ServiceID(serviceId),
		ID:             serviceapiapi.ServiceAPIID(serviceApiId),
		ProductTierID:  commonutils.ToPtr(serviceapiapi.ProductTierID(productTierId)),
		VersionSetName: versionSetName,
		VersionSetType: commonutils.ToPtr("Major"),
		IsPreferred:    isPreferred,
	}

	if err = serviceApi.ReleaseServiceAPI(context.Background(), request); err != nil {
		return
	}

	return
}
