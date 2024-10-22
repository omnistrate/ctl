package dataaccess

import (
	"context"

	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	serviceapiapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_apiapi"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/utils"
)

func DeleteProductTier(ctx context.Context, token, serviceID, productTierID string) (err error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	r, err := apiClient.ProductTierApiAPI.ProductTierApiDeleteProductTier(
		ctxWithToken,
		serviceID,
		productTierID,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return err
	}

	r.Body.Close()
	return nil
}

func DescribeProductTier(ctx context.Context, token, serviceID, productTierID string) (productTier *openapiclient.DescribeProductTierResult, err error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.ProductTierApiAPI.ProductTierApiDescribeProductTier(
		ctxWithToken,
		serviceID,
		productTierID,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return res, nil
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

	if err = serviceApi.ReleaseServiceAPI(ctx, request); err != nil {
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

	pendingChanges, err = serviceApi.DescribePendingChanges(ctx, request)
	if err != nil {
		return
	}

	return
}
