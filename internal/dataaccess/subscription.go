package dataaccess

import (
	"context"
	"github.com/pkg/errors"
	"net/http"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

func DescribeSubscription(ctx context.Context, token string, serviceID, environmentID, instanceID string) (resp *openapiclientfleet.FleetDescribeSubscriptionResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDescribeSubscription(
		ctxWithToken,
		serviceID,
		environmentID,
		instanceID,
	)

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	resp, r, err = req.Execute()
	if err != nil {
		return nil, handleFleetError(err)
	}
	return
}

func GetSubscriptionByCustomerEmail(ctx context.Context, token string, serviceID string, planID string, customerEmail string) (resp *openapiclientfleet.FleetDescribeSubscriptionResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	// Describe the service offering for this service and plan (product tier) ID to get the environment ID
	serviceOfferingResult, err := DescribeServiceOffering(ctx, token, serviceID, planID, "")
	if err != nil {
		return nil, handleFleetError(err)
	}

	var r *http.Response
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()

	for _, offering := range serviceOfferingResult.ConsumptionDescribeServiceOfferingResult.Offerings {
		if offering.ProductTierID == planID {
			req := apiClient.InventoryApiAPI.InventoryApiListSubscription(
				ctxWithToken,
				serviceID,
				offering.ServiceEnvironmentID,
			).ProductTierId(planID)

			var listSubscriptionResult *openapiclientfleet.FleetListSubscriptionsResult
			listSubscriptionResult, r, err = req.Execute()
			if err != nil {
				return nil, handleFleetError(err)
			}

			for _, subscription := range listSubscriptionResult.Subscriptions {
				if subscription.RootUserEmail == customerEmail {
					resp = &subscription
					return
				}
			}

			// Search user by email
			listUsersRes, _, err := apiClient.InventoryApiAPI.InventoryApiListAllUsers(ctxWithToken).Execute()
			if err != nil {
				return nil, handleFleetError(errors.Wrap(err, "failed to list users"))
			}

			userID := ""
			for _, user := range listUsersRes.Users {
				if *user.Email == customerEmail {
					userID = *user.UserId
					break
				}
			}

			if userID == "" {
				return nil, errors.Errorf("no user found with email %s", customerEmail)
			}

			// Subscription not found for the given customer email, create a new one
			createReq := apiClient.InventoryApiAPI.InventoryApiCreateSubscriptionOnBehalfOfCustomer(
				ctxWithToken,
				serviceID,
				offering.ServiceEnvironmentID,
			).FleetCreateSubscriptionOnBehalfOfCustomerRequest2(openapiclientfleet.FleetCreateSubscriptionOnBehalfOfCustomerRequest2{
				ProductTierId:            planID,
				OnBehalfOfCustomerUserId: userID,
			})

			createResp, _, err := createReq.Execute()
			if err != nil {
				return nil, handleFleetError(errors.Wrapf(err, "failed to create subscription for user %s", customerEmail))
			}

			// Describe the newly created subscription
			resp, err = DescribeSubscription(ctx, token, serviceID, offering.ServiceEnvironmentID, *createResp.Id)
			if err != nil {
				return nil, handleFleetError(errors.Wrapf(err, "failed to describe newly created subscription for user %s", customerEmail))
			}

			return resp, nil
		}
	}

	err = errors.New("no subscription found for the given customer email or the plan does not exist")
	return
}
