package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	producttierapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/product_tier_api"
	serviceapiapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_apiapi"
	tierversionsetapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/tier_version_set_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/utils"
)

func DeleteServicePlan(token, serviceId, productTierId string) (err error) {
	fleetService, err := httpclientwrapper.NewProductTier(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.DeleteProductTierRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceId),
		ID:        producttierapi.ProductTierID(productTierId),
	}

	if err = fleetService.DeleteProductTier(context.Background(), request); err != nil {
		return
	}

	return
}

func DescribeProductTier(token, serviceId, productTierId string) (productTier *producttierapi.DescribeProductTierResult, err error) {
	fleetService, err := httpclientwrapper.NewProductTier(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &producttierapi.DescribeProductTierRequest{
		Token:     token,
		ServiceID: producttierapi.ServiceID(serviceId),
		ID:        producttierapi.ProductTierID(productTierId),
	}

	productTier, err = fleetService.DescribeProductTier(context.Background(), request)
	if err != nil {
		return
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

func SetDefaultServicePlan(token, serviceId, productTierId, version string) (tierVersionSet *tierversionsetapi.TierVersionSet, err error) {
	versionSet, err := httpclientwrapper.NewVersionSet(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &tierversionsetapi.PromoteTierVersionSetRequest{
		Token:         token,
		ServiceID:     tierversionsetapi.ServiceID(serviceId),
		ProductTierID: tierversionsetapi.ProductTierID(productTierId),
		Version:       version,
	}

	if tierVersionSet, err = versionSet.PromoteTierVersionSet(context.Background(), request); err != nil {
		return
	}

	return
}
