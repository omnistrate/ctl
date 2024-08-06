package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	upgradepathapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	serviceapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/service_api"
	"github.com/omnistrate/ctl/utils"
)

func CreateUpgradePath(token, serviceID, productTierID, sourceVersion, targetVersion string, instanceIDs []string) (upgradepathapi.UpgradePathID, error) {
	upgradePath, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	res, err := upgradePath.CreateUpgradePath(context.Background(), &upgradepathapi.CreateUpgradePathRequest{
		Token:         token,
		ServiceID:     upgradepathapi.ServiceID(serviceID),
		ProductTierID: upgradepathapi.ProductTierID(productTierID),
		SourceVersion: sourceVersion,
		TargetVersion: targetVersion,
		UpgradeFilters: map[upgradepathapi.UpgradeFilterType][]string{
			"INSTANCE_IDS": instanceIDs,
		},
	})
	if err != nil {
		return "", err
	}

	return res.UpgradePathID, nil
}

func DescribeUpgradePath(token, serviceID, productTierID, upgradePathID string) (*upgradepathapi.UpgradePath, error) {
	upgradePath, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := upgradePath.DescribeUpgradePath(context.Background(), &upgradepathapi.DescribeUpgradePathRequest{
		Token:         token,
		ServiceID:     upgradepathapi.ServiceID(serviceID),
		ProductTierID: upgradepathapi.ProductTierID(productTierID),
		UpgradePathID: upgradepathapi.UpgradePathID(upgradePathID),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func ListUpgradePaths(token, serviceID, productTierID string) ([]*upgradepathapi.UpgradePath, error) {
	upgradePath, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := upgradePath.ListUpgradePath(context.Background(), &upgradepathapi.ListUpgradePathsRequest{
		Token:         token,
		ServiceID:     upgradepathapi.ServiceID(serviceID),
		ProductTierID: upgradepathapi.ProductTierID(productTierID),
	})
	if err != nil {
		return nil, err
	}

	return res.UpgradePaths, nil
}

func ListAllUpgradePaths(token string) ([]*upgradepathapi.UpgradePath, error) {
	svc, err := httpclientwrapper.NewService(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := svc.ListService(context.Background(), &serviceapi.List{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	upgradePath, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	allUpgradePaths := make([]*upgradepathapi.UpgradePath, 0)

	for _, service := range res.Services {
		for _, env := range service.ServiceEnvironments {
			for _, tier := range env.ServicePlans {
				upgradePaths, err := upgradePath.ListUpgradePath(context.Background(), &upgradepathapi.ListUpgradePathsRequest{
					Token:         token,
					ServiceID:     upgradepathapi.ServiceID(service.ID),
					ProductTierID: upgradepathapi.ProductTierID(tier.ProductTierID),
				})
				if err != nil {
					return nil, err
				}

				allUpgradePaths = append(allUpgradePaths, upgradePaths.UpgradePaths...)
			}
		}
	}

	return allUpgradePaths, nil
}

func ListEligibleInstancesPerUpgrade(token, serviceID, productTierID, upgradePathID string) ([]*upgradepathapi.InstanceUpgrade, error) {
	upgradePath, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := upgradePath.ListEligibleInstancesPerUpgrade(context.Background(), &upgradepathapi.ListEligibleInstancesPerUpgradeRequest{
		Token:         token,
		ServiceID:     upgradepathapi.ServiceID(serviceID),
		ProductTierID: upgradepathapi.ProductTierID(productTierID),
		UpgradePathID: upgradepathapi.UpgradePathID(upgradePathID),
	})
	if err != nil {
		return nil, err
	}

	return res.Instances, nil
}
