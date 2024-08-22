package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	resourceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/resource_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeResource(token, serviceId, resourceId string, productTierId, productTierVersion *string) (resource *resourceapi.DescribeResourceResult, err error) {
	service, err := httpclientwrapper.NewResource(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return
	}

	request := &resourceapi.DescribeResourceRequest{
		Token:              token,
		ServiceID:          resourceapi.ServiceID(serviceId),
		ID:                 resourceapi.ResourceID(resourceId),
		ProductTierVersion: productTierVersion,
		ProductTierID:      (*resourceapi.ProductTierID)(productTierId),
	}

	if resource, err = service.DescribeResource(context.Background(), request); err != nil {
		return
	}

	return
}
