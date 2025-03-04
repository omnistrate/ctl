package dataaccess

import (
	"context"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	upgradepathapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/internal/config"
)

func CreateUpgradePath(ctx context.Context, token, serviceID, productTierID, sourceVersion, targetVersion string, scheduledDate *string, instanceIDs []string) (upgradepathapi.UpgradePathID, error) {
	upgradePath, err := httpclientwrapper.NewInventory(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return "", err
	}

	res, err := upgradePath.CreateUpgradePath(ctx, &upgradepathapi.CreateUpgradePathRequest{
		Token:         token,
		ServiceID:     upgradepathapi.ServiceID(serviceID),
		ProductTierID: upgradepathapi.ProductTierID(productTierID),
		SourceVersion: sourceVersion,
		TargetVersion: targetVersion,
		ScheduledDate: scheduledDate,
		UpgradeFilters: map[upgradepathapi.UpgradeFilterType][]string{
			"INSTANCE_IDS": instanceIDs,
		},
	})
	if err != nil {
		return "", err
	}

	return res.UpgradePathID, nil
}

func DescribeUpgradePath(ctx context.Context, token, serviceID, productTierID, upgradePathID string) (*openapiclientfleet.UpgradePath, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiDescribeUpgradePath(
		ctxWithToken,
		serviceID,
		productTierID,
		upgradePathID,
	)

	resp, r, err := req.Execute()
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return resp, nil
}

func ListEligibleInstancesPerUpgrade(ctx context.Context, token, serviceID, productTierID, upgradePathID string) ([]openapiclientfleet.InstanceUpgrade, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiListEligibleInstancesPerUpgrade(
		ctxWithToken,
		serviceID,
		productTierID,
		upgradePathID,
	)

	resp, r, err := req.Execute()
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()
	if err != nil {
		return nil, handleFleetError(err)
	}

	return resp.GetInstances(), nil
}
