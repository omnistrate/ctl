package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	tierversionsetapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/tier_version_set_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/pkg/errors"
)

func FindLatestVersion(ctx context.Context, token, serviceID, productTierID string) (string, error) {
	versionSet, err := httpclientwrapper.NewVersionSet(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return "", err
	}

	res, err := versionSet.ListTierVersionSets(ctx, &tierversionsetapi.ListTierVersionSetsRequest{
		Token:                  token,
		ServiceID:              tierversionsetapi.ServiceID(serviceID),
		ProductTierID:          tierversionsetapi.ProductTierID(productTierID),
		LatestMajorVersionOnly: utils.ToPtr(true),
	})
	if err != nil {
		return "", err
	}

	if len(res.TierVersionSets) == 0 {
		return "", errors.New("no version found")
	}

	return res.TierVersionSets[0].Version, nil
}

func FindPreferredVersion(ctx context.Context, token, serviceID, productTierID string) (string, error) {
	versionSet, err := httpclientwrapper.NewVersionSet(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return "", err
	}

	res, err := versionSet.ListTierVersionSets(ctx, &tierversionsetapi.ListTierVersionSetsRequest{
		Token:         token,
		ServiceID:     tierversionsetapi.ServiceID(serviceID),
		ProductTierID: tierversionsetapi.ProductTierID(productTierID),
	})
	if err != nil {
		return "", err
	}

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

func DescribeVersionSet(ctx context.Context, token, serviceID, productTierID, version string) (*tierversionsetapi.TierVersionSet, error) {
	versionSet, err := httpclientwrapper.NewVersionSet(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := versionSet.DescribeTierVersionSet(ctx, &tierversionsetapi.DescribeTierVersionSetRequest{
		Token:         token,
		ServiceID:     tierversionsetapi.ServiceID(serviceID),
		ProductTierID: tierversionsetapi.ProductTierID(productTierID),
		Version:       version,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func SetDefaultServicePlan(ctx context.Context, token, serviceID, productTierID, version string) (tierVersionSet *tierversionsetapi.TierVersionSet, err error) {
	versionSet, err := httpclientwrapper.NewVersionSet(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return
	}

	request := &tierversionsetapi.PromoteTierVersionSetRequest{
		Token:         token,
		ServiceID:     tierversionsetapi.ServiceID(serviceID),
		ProductTierID: tierversionsetapi.ProductTierID(productTierID),
		Version:       version,
	}

	if tierVersionSet, err = versionSet.PromoteTierVersionSet(ctx, request); err != nil {
		return
	}

	return
}
