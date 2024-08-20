package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeSubscription(token string, serviceID, environmentID, instanceID string) (*inventoryapi.FleetDescribeSubscriptionResult, error) {
	subscription, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := inventoryapi.FleetDescribeSubscriptionRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		ID:            inventoryapi.SubscriptionID(instanceID),
	}

	res, err := subscription.DescribeSubscription(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
