package cmd

import (
	"fmt"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_build_basic(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	// Step 1: login
	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	// Step 2: get compose files
	composeFiles, err := os.ReadDir("../composefiles")
	require.NoError(err)

	// Step 3: test build service on all compose files
	for _, f := range composeFiles {
		if f.IsDir() {
			continue
		}

		rootCmd.SetArgs([]string{"build", "-f", "../composefiles/" + f.Name(), "--name", f.Name(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
		err = rootCmd.Execute()
		require.NoError(err, f.Name())

		rootCmd.SetArgs([]string{"describe", "--service-id", serviceID})
		err = rootCmd.Execute()
		require.NoError(err, f.Name())

		rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
		err = rootCmd.Execute()
		require.NoError(err, f.Name())
	}
}

func Test_build_invalid_file(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "invalid_file.yaml", "--name", "My Service", "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "file does not exist: invalid_file.yaml")
}

func Test_build_no_file(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "--name", "My Service", "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "must provide --file or -f")
}

func Test_build_invalid_file_format(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "invalid_file.txt", "--name", "My Service", "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "file must be a valid docker-compose file in .yaml or .yml format")
}

func Test_build_create_no_name(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "name is required")
}

func Test_build_create_no_description(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--name", "cassandra", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "Invalid request: description: parameter is empty")
}

func Test_build_create_no_service_logo_url(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/cassandra.yaml", "--name", "cassandra", "--description", "My Service Description"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}
