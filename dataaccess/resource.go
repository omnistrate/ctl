package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	resourceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/resource_api"
	"github.com/omnistrate/ctl/config"
)

func DescribeResource(token, serviceID, resourceID string, productTierID, productTierVersion *string) (resource *resourceapi.DescribeResourceResult, err error) {
	service, err := httpclientwrapper.NewResource(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &resourceapi.DescribeResourceRequest{
		Token:              token,
		ServiceID:          resourceapi.ServiceID(serviceID),
		ID:                 resourceapi.ResourceID(resourceID),
		ProductTierVersion: productTierVersion,
		ProductTierID:      (*resourceapi.ProductTierID)(productTierID),
	}

	if resource, err = service.DescribeResource(context.Background(), request); err != nil {
		return
	}

	return
}
