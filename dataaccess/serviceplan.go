package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	producttierapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/product_tier_api"
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
