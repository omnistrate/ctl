package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	tierversionsetapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/tier_version_set_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
)

func FindLatestVersion(token, serviceID, productTierID string) (string, error) {
	versionSet, err := httpclientwrapper.NewVersionSet(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	res, err := versionSet.ListTierVersionSets(context.Background(), &tierversionsetapi.ListTierVersionSetsRequest{
		Token:                  token,
		ServiceID:              tierversionsetapi.ServiceID(serviceID),
		ProductTierID:          tierversionsetapi.ProductTierID(productTierID),
		LatestMajorVersionOnly: commonutils.ToPtr(true),
	})
	if err != nil {
		return "", err
	}

	if len(res.TierVersionSets) == 0 {
		return "", errors.New("no version found")
	}

	return res.TierVersionSets[0].Version, nil
}

func FindPreferredVersion(token, serviceID, productTierID string) (string, error) {
	versionSet, err := httpclientwrapper.NewVersionSet(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	res, err := versionSet.ListTierVersionSets(context.Background(), &tierversionsetapi.ListTierVersionSetsRequest{
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

func DescribeVersionSet(token, serviceID, productTierID, version string) (*tierversionsetapi.TierVersionSet, error) {
	versionSet, err := httpclientwrapper.NewVersionSet(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := versionSet.DescribeTierVersionSet(context.Background(), &tierversionsetapi.DescribeTierVersionSetRequest{
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
