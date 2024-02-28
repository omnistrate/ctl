package cmd

import (
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_describe_basic(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/ferretdb.yaml", "--name", "ferretdb", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"describe"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove"})
	err = rootCmd.Execute()
	require.NoError(err)
}

func Test_describe_no_service(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"describe"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "service does not exist")
}
