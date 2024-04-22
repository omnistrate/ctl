package cmd

import (
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_remove_basic(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=Test@1234"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--name", "cassandra", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "must provide --service-id")
}
