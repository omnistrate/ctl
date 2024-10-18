package customnetwork

import (
	"context"
	"fmt"
	"testing"

	"github.com/omnistrate/ctl/cmd"
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

	// PASS: list custom networks
	cmd.RootCmd.SetArgs([]string{"custom-network", "list", "--filter", "cloud_provider:aws,region:ap-south-1"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)
}
