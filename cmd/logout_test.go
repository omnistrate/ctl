package cmd

import (
	"fmt"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_logout(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	// FAIL: logout without login
	rootCmd.SetArgs([]string{"logout"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "config file not found")

	// PASS: logout after login
	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"logout"})
	err = rootCmd.Execute()
	require.NoError(err)
}
