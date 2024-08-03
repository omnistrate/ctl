package dataaccess

import (
	"bytes"
	"encoding/json"
	"fmt"
	upgradepathapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
	"net/http"
)

func CreateUpgradePath(token, serviceID, productTierID, sourceVersion, TargetVersion, instanceID string) (upgradepathapi.UpgradePathID, error) {
	url := fmt.Sprintf("%s://%s/2022-09-01-00/fleet/service/%v/productTier/%v/upgrade-path", utils.GetHostScheme(), utils.GetHost(), serviceID, productTierID)

	payload := upgradepathapi.CreateUpgradePathRequest{
		Token:         token,
		ServiceID:     upgradepathapi.ServiceID(serviceID),
		ProductTierID: upgradepathapi.ProductTierID(productTierID),
		SourceVersion: sourceVersion,
		TargetVersion: TargetVersion,
		UpgradeFilters: map[upgradepathapi.UpgradeFilterType][]string{
			"INSTANCE_IDS": {
				instanceID,
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%s", resp.Status)
		return "", err
	}

	// Marshal the response body into a struct
	var result upgradepathapi.UpgradePath
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.UpgradePathID, nil
}

func DescribeUpgradePath(token, serviceID, productTierID, upgradePathID string) (*upgradepathapi.UpgradePath, error) {
	url := fmt.Sprintf("%s://%s/2022-09-01-00/fleet/service/%v/productTier/%v/upgrade-path/%v", utils.GetHostScheme(), utils.GetHost(), serviceID, productTierID, upgradePathID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("%s", resp.Status)
		return nil, err
	}

	// Marshal the response body into a struct
	var result upgradepathapi.UpgradePath
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
