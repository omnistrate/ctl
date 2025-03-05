package dataaccess

import (
	"context"
	"github.com/omnistrate/ctl/internal/model"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

func CreateUpgradePath(ctx context.Context, token, serviceID, productTierID, sourceVersion, targetVersion string, scheduledDate *string, instanceIDs []string) (string, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiCreateUpgradePath(ctxWithToken, serviceID, productTierID).
		CreateUpgradePathRequest2(openapiclientfleet.CreateUpgradePathRequest2{
			SourceVersion: sourceVersion,
			TargetVersion: targetVersion,
			ScheduledDate: scheduledDate,
			UpgradeFilters: map[string][]string{
				"INSTANCE_IDS": instanceIDs,
			},
		})

	resp, r, err := req.Execute()
	defer func() {
		if r != nil {
			_ = r.Body.Close()
		}
	}()
	if err != nil {
		return "", handleFleetError(err)
	}

	return resp.UpgradePathId, nil
}

func ManageLifecycle(ctx context.Context, token, serviceID, productTierID, upgradePathID string, action model.UpgradeMaintenanceAction) (*openapiclientfleet.UpgradePath, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)
	apiClient := getFleetClient()

	req := apiClient.InventoryApiAPI.InventoryApiManageUpgradePath(
		ctxWithToken,
		serviceID,
		productTierID,
		upgradePathID,
	)
	req = req.ManageUpgradePathLifecycleRequest2(openapiclientfleet.ManageUpgradePathLifecycleRequest2{
		Action: action.String(),
	})
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
