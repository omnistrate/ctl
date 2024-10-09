package customnetwork

import (
	"context"
	"fmt"
	"testing"

	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/customnetwork"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/omnistrate/ctl/test/testutils"
	"github.com/stretchr/testify/require"
)

func Test_custom_network_lifecycle(t *testing.T) {
	testutils.SmokeTest(t)

	ctx := context.TODO()

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// Pre-test cleanup
	token, err := config.GetToken()
	require.NoError(err)
	deleteCustomNetworkIfExists(t, ctx, token, "aws", "ap-south-1", "ctl-test-network")

	// PASS: create custom network
	cmd.RootCmd.SetArgs([]string{"custom-network", "create", "--cloud-provider", "aws", "--region", "ap-south-1", "--cidr", "1.2.255.1/16", "--name", "ctl-test-network"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	customNetworkID := customnetwork.CustomNetworkID
	require.NotEmpty(customNetworkID)

	// PASS: describe custom network
	cmd.RootCmd.SetArgs([]string{"custom-network", "describe", "--custom-network-id", customNetworkID})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: describe custom network by name
	cmd.RootCmd.SetArgs([]string{"custom-network", "describe", "ctl-test-network"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: describe manually
	customNetwork, err := dataaccess.DescribeCustomNetwork(ctx, token, customnetworkapi.DescribeCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(customNetworkID),
	})

	require.NoError(err)
	require.NotNil(customNetwork)
	require.Equal(customNetworkID, string(customNetwork.ID))
	require.NotNil(customNetwork.Name)
	require.Equal("ctl-test-network", *customNetwork.Name)
	require.Equal("1.2.255.1/16", customNetwork.Cidr)
	require.Equal("aws", string(customNetwork.CloudProviderName))
	require.Equal("ap-south-1", customNetwork.CloudProviderRegion)

	// PASS: list custom networks
	cmd.RootCmd.SetArgs([]string{"custom-network", "list", "--filter", "cloud_provider:aws,region:ap-south-1"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: delete custom network by name
	cmd.RootCmd.SetArgs([]string{"custom-network", "delete", "ctl-test-network"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// FAIL: describe again
	customNetwork, err = dataaccess.DescribeCustomNetwork(ctx, token, customnetworkapi.DescribeCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(customNetworkID),
	})
	require.Error(err)
	require.Nil(customNetwork)
}

func deleteCustomNetworkIfExists(t *testing.T, ctx context.Context, token, cloudProvider, region, customNetworkName string) {
	customNetworks, err := dataaccess.ListCustomNetworks(ctx, token, customnetworkapi.ListCustomNetworksRequest{
		CloudProviderName:   utils.ToPtr(customnetworkapi.CloudProvider(cloudProvider)),
		CloudProviderRegion: utils.ToPtr(region),
	})
	require.NoError(t, err)

	for _, network := range customNetworks.CustomNetworks {
		if network.Name != nil && *network.Name == customNetworkName {
			err = dataaccess.DeleteCustomNetwork(ctx, token, customnetworkapi.DeleteCustomNetworkRequest{
				ID: network.ID,
			})
			require.NoError(t, err)
		}
	}
}
