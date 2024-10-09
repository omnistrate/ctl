package upgrade

import (
	"fmt"
	"testing"

	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/test/testutils"

	"github.com/stretchr/testify/require"
)

func Test_upgrade_basic(t *testing.T) {
	testutils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	testEmail, testPassword, err := testutils.GetTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"upgrade", "instance-invalid", "--version", "latest"})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "instance-invalid not found. Please check the instance ID and try again")

	cmd.RootCmd.SetArgs([]string{"upgrade", "status", "upgrade-invalid"})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "upgrade-invalid not found")

	// TODO: Create real upgrade path after we added the CRUD instance cmd
	cmd.RootCmd.SetArgs([]string{"upgrade", "status", "upgrade-qtxOTgcnDI"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"upgrade", "status", "upgrade-qtxOTgcnDI", "--output", "json"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"upgrade", "status", "detail", "upgrade-qtxOTgcnDI"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"upgrade", "status", "detail", "upgrade-qtxOTgcnDI", "--output", "json"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}
