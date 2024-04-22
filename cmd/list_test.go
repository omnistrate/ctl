package cmd

import (
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_list_basic(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=Test@1234"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/ferretdb.yaml", "--name", "ferretdb", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"list"})
	err = rootCmd.Execute()
	require.NoError(err)
}

func Test_list_no_service(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=Test@1234"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"list"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "service does not exist")
}

func Test_list_no_service_logo_url(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+cli@omnistrate.com", "--password=Test@1234"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/ferretdb.yaml", "--name", "ferretdb", "--description", "My Service Description"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"list"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove"})
	err = rootCmd.Execute()
	require.NoError(err)
}
