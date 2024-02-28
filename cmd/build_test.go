package cmd

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_build(t *testing.T) {
	require := require.New(t)
	defer cleanup()

	// Step 1: login
	rootCmd.SetArgs([]string{"login", "--email=xzhang+ctltest@omnistrate.com", "--password=ctltest"})
	err := rootCmd.Execute()
	require.NoError(err)

	// Step 2: get compose files
	composeFiles, err := os.ReadDir("../composefiles")
	require.NoError(err)

	// Step 3: test build service on all compose files
	for _, f := range composeFiles {
		if f.Name() == "mariadbcluster.yaml" || f.Name() == "postgres_advanced_serverless.yaml" {
			continue
		}

		rootCmd.SetArgs([]string{"build", "-f", "../composefiles/" + f.Name(), "--name", f.Name(), "--description", "My Service Description"})
		err = rootCmd.Execute()
		require.NoError(err)

		rootCmd.SetArgs([]string{"remove"})
		err = rootCmd.Execute()
		require.NoError(err)
	}
}
