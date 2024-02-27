package cmd

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_logout(t *testing.T) {
	require := require.New(t)

	// FAIL: logout without login
	rootCmd.SetArgs([]string{"logout"})
	err := rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "no auth config found")

	// PASS: logout after login
	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"logout"})
	err = rootCmd.Execute()
	require.NoError(err)
}
