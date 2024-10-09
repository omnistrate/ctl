package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	upgradepathapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
)

func CreateUpgradePath(token, serviceID, productTierID, sourceVersion, targetVersion string, instanceIDs []string) (upgradepathapi.UpgradePathID, error) {
	upgradePath, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
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
	upgradePath, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
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

func ListEligibleInstancesPerUpgrade(token, serviceID, productTierID, upgradePathID string) ([]*upgradepathapi.InstanceUpgrade, error) {
	upgradePath, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
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
