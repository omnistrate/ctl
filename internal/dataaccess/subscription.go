package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
)

func DescribeSubscription(ctx context.Context, token string, serviceID, environmentID, instanceID string) (*inventoryapi.FleetDescribeSubscriptionResult, error) {
	subscription, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	request := inventoryapi.FleetDescribeSubscriptionRequest{
		Token:         token,
		ServiceID:     inventoryapi.ServiceID(serviceID),
		EnvironmentID: inventoryapi.ServiceEnvironmentID(environmentID),
		ID:            inventoryapi.SubscriptionID(instanceID),
	}

	res, err := subscription.DescribeSubscription(ctx, &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
