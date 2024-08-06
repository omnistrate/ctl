package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	inventoryapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/inventory_api"
	"github.com/omnistrate/ctl/utils"
)

func SearchInventory(token, query string) (*inventoryapi.SearchInventoryResult, error) {
	inventory, err := httpclientwrapper.NewInventory(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := inventory.SearchInventory(context.Background(), &inventoryapi.SearchInventoryRequest{
		Token: token,
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
