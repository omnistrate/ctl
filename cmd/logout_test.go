package cmd

import (
	"fmt"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
)

func Test_logout(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	// FAIL: logout without login
	RootCmd.SetArgs([]string{"logout"})
	err = RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "config file not found")

	// PASS: logout after login
	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)

	RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = RootCmd.Execute()
	require.NoError(err)

	RootCmd.SetArgs([]string{"logout"})
	err = RootCmd.Execute()
	require.NoError(err)
}
