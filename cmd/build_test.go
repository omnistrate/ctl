package cmd

import (
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_build(t *testing.T) {
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
