package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
)

func Test_build_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	// Step 1: login
	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
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

		rootCmd.SetArgs([]string{"build", "-f", "../composefiles/" + f.Name(), "--name", f.Name() + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
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

func Test_build_update_service(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	// PASS: login
	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: create mysql cluster service
	serviceName := "mysql cluster" + uuid.NewString()
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png", "--release"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"describe", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: update mysql cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_variation_apiparam_image_infra_capability.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: update back to original mysql cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: update mysql cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_variation_account_integration_resource.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: update back to original mysql cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: create postgres cluster service
	serviceName = "postgres cluster" + uuid.NewString()
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/postgrescluster_original.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png", "--release"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"describe", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: update postgres cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/postgrescluster_variation_load_balancer.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: update back to original postgres cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/postgrescluster_original.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	// PASS: add new service plan to postgres cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/postgrescluster_variation_new_tier.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}

func Test_build_duplicate_service_plan_name(t *testing.T) {
	t.Skip("Skipping until we fix this logic")

	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	// PASS: login
	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	serviceName := "mysql cluster" + uuid.NewString()
	// PASS: create mysql cluster service in dev environment
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png", "--release"})
	err = rootCmd.Execute()
	require.NoError(err)
	devEnvironmentID := environmentID
	devProductTierID := productTierID
	require.NotEmpty(devEnvironmentID)
	require.NotEmpty(devProductTierID)

	// PASS: create mysql cluster service in prod environment
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--environment", "prod", "--release"})
	err = rootCmd.Execute()
	require.NoError(err)
	prodEnvironmentID := environmentID
	prodProductTierID := productTierID
	require.NotEmpty(prodEnvironmentID)
	require.NotEmpty(prodEnvironmentID)
	require.NotEqual(devEnvironmentID, prodEnvironmentID)
	require.NotEqual(devProductTierID, prodProductTierID)

	// PASS: update dev mysql cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_variation_apiparam_image_infra_capability.yaml", "--name", serviceName, "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)
	require.Equal(devEnvironmentID, environmentID)
	require.Equal(devProductTierID, productTierID)

	// PASS: update prod mysql cluster service
	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/variations/mysqlcluster_variation_apiparam_image_infra_capability.yaml", "--name", serviceName, "--environment", "prod", "--release-as-preferred"})
	err = rootCmd.Execute()
	require.NoError(err)
	require.Equal(prodEnvironmentID, environmentID)
	require.Equal(prodProductTierID, productTierID)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}

func Test_build_invalid_file(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "invalid_file.yaml", "--name", "My Service" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "no such file or directory")
}

func Test_build_no_file(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "--name", "My Service" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "must provide --file or -f")
}

func Test_build_create_no_name(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/postgresql.yaml", "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "name is required")
}

func Test_build_create_no_description(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}

func Test_build_create_no_service_logo_url(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--description", "My Service Description"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}
