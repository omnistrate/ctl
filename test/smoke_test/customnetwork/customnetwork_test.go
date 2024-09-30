package customnetwork

import (
	"fmt"
	customnetworkapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/custom_network_api"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/customnetwork"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/test/testutils"
	ctlutils "github.com/omnistrate/ctl/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_custom_network_lifecycle(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Pre-test cleanup
	token, err := ctlutils.GetToken()
	require.NoError(err)
	deleteCustomNetworkIfExists(t, token, "aws", "ap-south-1", "ctl-test-network")

	// PASS: create custom network
	cmd.RootCmd.SetArgs([]string{"custom-network", "create", "--cloud-provider", "aws", "--region", "ap-south-1", "--cidr", "10.99.101.0/24", "--name", "ctl-test-network"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	customNetworkID := customnetwork.CustomNetworkID
	require.NotEmpty(customNetworkID)

	// PASS: describe custom network
	cmd.RootCmd.SetArgs([]string{"custom-network", "describe", customNetworkID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe manually
	customNetwork, err := dataaccess.DescribeCustomNetwork(token, customnetworkapi.DescribeCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(customNetworkID),
	})

	require.NoError(err)
	require.NotNil(customNetwork)
	require.Equal(customNetworkID, string(customNetwork.ID))
	require.NotNil(customNetwork.Name)
	require.Equal("ctl-test-network", *customNetwork.Name)
	require.Equal("10.99.101.0/24", customNetwork.Cidr)
	require.Equal("aws", string(customNetwork.CloudProviderName))
	require.Equal("ap-south-1", customNetwork.CloudProviderRegion)

	// PASS: list custom networks
	cmd.RootCmd.SetArgs([]string{"custom-network", "list", "--cloud-provider", "aws", "--region", "ap-south-1"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: delete custom network
	cmd.RootCmd.SetArgs([]string{"custom-network", "delete", customNetworkID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// FAIL: describe again
	customNetwork, err = dataaccess.DescribeCustomNetwork(token, customnetworkapi.DescribeCustomNetworkRequest{
		ID: customnetworkapi.CustomNetworkID(customNetworkID),
	})
	require.Error(err)
	require.Nil(customNetwork)
}

func deleteCustomNetworkIfExists(t *testing.T, token, cloudProvider, region, customNetworkName string) {
	customNetworks, err := dataaccess.ListCustomNetworks(token, customnetworkapi.ListCustomNetworksRequest{
		CloudProviderName:   utils.ToPtr(customnetworkapi.CloudProvider(cloudProvider)),
		CloudProviderRegion: utils.ToPtr(region),
	})
	require.NoError(t, err)

	for _, network := range customNetworks.CustomNetworks {
		if network.Name != nil && *network.Name == customNetworkName {
			err = dataaccess.DeleteCustomNetwork(token, customnetworkapi.DeleteCustomNetworkRequest{
				ID: network.ID,
			})
			require.NoError(t, err)
		}
	}
}
