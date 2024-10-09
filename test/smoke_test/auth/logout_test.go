package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/omnistrate/ctl/config"

	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/test/testutils"

	"github.com/stretchr/testify/require"
)

func Test_logout(t *testing.T) {
	testutils.SmokeTest(t)

	ctx := context.TODO()

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	// FAIL: logout without login
	cmd.RootCmd.SetArgs([]string{"logout"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), config.ErrConfigFileNotFound.Error())

	// PASS: logout after login
	testEmail, testPassword, err := testutils.GetTestAccount()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"logout"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)
}
