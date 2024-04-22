package cmd

import (
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_logout(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	// FAIL: logout without login
	rootCmd.SetArgs([]string{"logout"})
	err := rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "config file not found")

	// PASS: logout after login
	rootCmd.SetArgs([]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=Test@1234"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"logout"})
	err = rootCmd.Execute()
	require.NoError(err)
}
