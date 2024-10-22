package dataaccess

import (
	"context"

	openapiclientfleet "github.com/omnistrate-oss/omnistrate-sdk-go/fleet"
)

func SearchInventory(ctx context.Context, token, query string) (*openapiclientfleet.SearchInventoryResult, error) {
	ctxWithToken := context.WithValue(ctx, openapiclientfleet.ContextAccessToken, token)

	req := *openapiclientfleet.NewSearchServiceInventoryRequestBody(query)

	apiClient := getFleetClient()
	res, r, err := apiClient.InventoryApiAPI.
		InventoryApiSearchInventory(ctxWithToken).
		SearchServiceInventoryRequestBody(req).
		Execute()

	err = handleFleetError(err)
	if err != nil {
		return nil, err
	}

	r.Body.Close()
	return res, nil
}
