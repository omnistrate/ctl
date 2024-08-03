package dataaccess

import (
	"bytes"
	"encoding/json"
	"fmt"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
	"net/http"
)

func SearchInventory(token, query string) (*inventoryapi.SearchInventoryResult, error) {
	url := fmt.Sprintf("%s://%s/2022-09-01-00/fleet/search-inventory", utils.GetHostScheme(), utils.GetHost())

	payload := inventoryapi.SearchInventoryRequest{
		Query: query,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
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
	var result inventoryapi.SearchInventoryResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
