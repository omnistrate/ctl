package dataaccess

import (
	"context"

	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/pkg/errors"
)

func FindLatestVersion(ctx context.Context, token, serviceID, productTierID string) (string, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.TierVersionSetApiAPI.TierVersionSetApiListTierVersionSets(
		ctxWithToken,
		serviceID,
		productTierID,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	if len(res.TierVersionSets) == 0 {
		return "", errors.New("no version found")
	}

	return res.TierVersionSets[0].Version, nil
}

func FindPreferredVersion(ctx context.Context, token, serviceID, productTierID string) (string, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.TierVersionSetApiAPI.TierVersionSetApiListTierVersionSets(
		ctxWithToken,
		serviceID,
		productTierID,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	defer r.Body.Close()

	if len(res.TierVersionSets) == 0 {
		return "", errors.New("no version found")
	}

	for _, versionSet := range res.TierVersionSets {
		if versionSet.Status == "Preferred" {
			return versionSet.Version, nil
		}
	}

	return "", errors.New("no preferred version found")
}

func DescribeVersionSet(ctx context.Context, token, serviceID, productTierID, version string) (*openapiclient.TierVersionSet, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.TierVersionSetApiAPI.TierVersionSetApiDescribeTierVersionSet(
		ctxWithToken,
		serviceID,
		productTierID,
		version,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	return res, nil
}

func SetDefaultServicePlan(ctx context.Context, token, serviceID, productTierID, version string) (*openapiclient.TierVersionSet, error) {
	ctxWithToken := context.WithValue(ctx, openapiclient.ContextAccessToken, token)

	apiClient := getV1Client()
	res, r, err := apiClient.TierVersionSetApiAPI.TierVersionSetApiPromoteTierVersionSet(
		ctxWithToken,
		serviceID,
		productTierID,
		version,
	).Execute()

	err = handleV1Error(err)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	return res, nil
}
