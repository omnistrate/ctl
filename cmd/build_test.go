package cmd

import (
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_build_basic(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	// Step 1: login
	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	// Step 2: get compose files
	composeFiles, err := os.ReadDir("../composefiles")
	require.NoError(err)

	// Step 3: test build service on all compose files
	skipComposeFiles := []string{"mariadbcluster.yaml", "postgres_advanced_serverless.yaml", "opensearch.yaml", "postgrescluster.yaml"}
	for _, f := range composeFiles {
		if testutils.Contains(skipComposeFiles, f.Name()) {
			continue
		}

		rootCmd.SetArgs([]string{"build", "-f", "../composefiles/" + f.Name(), "--name", f.Name(), "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
		err = rootCmd.Execute()
		require.NoError(err)

		rootCmd.SetArgs([]string{"describe"})
		err = rootCmd.Execute()
		require.NoError(err)

		rootCmd.SetArgs([]string{"remove"})
		err = rootCmd.Execute()
		require.NoError(err)
	}
}

func Test_build_invalid_file(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "invalid_file.yaml", "--name", "My Service", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "file does not exist: invalid_file.yaml")
}

func Test_build_no_file(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "--name", "My Service", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "must provide --file or -f")
}

func Test_build_invalid_file_format(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "invalid_file.txt", "--name", "My Service", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "file must be a valid docker-compose file in .yaml or .yml format")
}

func Test_build_create_no_name(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--description", "My Service Description", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "name is required for creating service")
}

func Test_build_create_no_description(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--name", "cassandra", "--service-logo-url", "https://my-service-logo.com/logo.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "description is required for creating service")
}

// TODO: fix this test
//func Test_build_create_no_service_logo_url(t *testing.T) {
//	require := require.New(t)
//	defer testutils.Cleanup()
//
//	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
//	err := rootCmd.Execute()
//	require.NoError(err)
//
//	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--name", "cassandra", "--description", "My Service Description"})
//	err = rootCmd.Execute()
//	require.NoError(err)
//}
