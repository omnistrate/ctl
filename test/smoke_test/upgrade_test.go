package smoke

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/test/testutils"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_upgrade_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	//testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	//require.NoError(err)
	//cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	//err = cmd.RootCmd.Execute()
	//require.NoError(err)

	//cmd.RootCmd.SetArgs([]string{"upgrade", "instance-7yx75x28n", "--version", "latest"})
	//err = cmd.RootCmd.Execute()
	//require.NoError(err)

	_, err := dataaccess.SearchInventory("token", "query")
	require.Error(err)

	_, err = dataaccess.DescribeResourceInstance("token", "serviceID", "environmentID", "instanceID")
	require.Error(err)

	_, err = dataaccess.CreateUpgradePath("token", "serviceID", "productTierID", "sourceVersion", "TargetVersion", "instanceID")
	require.Error(err)

	_, err = dataaccess.ListUpgradePaths("token", "serviceID", "productTierID")
	require.Error(err)

	_, err = dataaccess.DescribeUpgradePath("token", "serviceID", "productTierID", "upgradePathID")
	require.Error(err)

	_, err = dataaccess.FindLatestVersion("token", "serviceID", "productTierID")
	require.Error(err)

	_, err = dataaccess.DescribeVersionSet("token", "serviceID", "productTierID", "version")
	require.Error(err)

	_, err = dataaccess.ListDomains("token")
	require.Error(err)
}
