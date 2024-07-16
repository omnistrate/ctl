package smoke

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/test/testutils"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/omnistrate/commons/pkg/utils"
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
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Step 2: get compose files
	composeFiles, err := os.ReadDir("composefiles")
	require.NoError(err)

	// Step 3: test build service on all compose files
	for _, f := range composeFiles {
		if f.IsDir() {
			continue
		}

		cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/" + f.Name(), "--name", f.Name() + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
		err = cmd.RootCmd.Execute()
		require.NoError(err, f.Name())

		cmd.RootCmd.SetArgs([]string{"describe", "service", f.Name()})
		err = cmd.RootCmd.Execute()
		require.NoError(err, f.Name())

		cmd.RootCmd.SetArgs([]string{"remove", "--service-id", cmd.ServiceID})
		err = cmd.RootCmd.Execute()
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
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: create mysql cluster service
	serviceName := "mysql cluster" + uuid.NewString()
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png", "--release"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"describe", "service", serviceName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: update mysql cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_variation_apiparam_image_infra_capability.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: update back to original mysql cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: update mysql cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_variation_account_integration_resource.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: update back to original mysql cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"remove", "--service-id", cmd.ServiceID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: create postgres cluster service
	serviceName = "postgres cluster" + uuid.NewString()
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/postgrescluster_original.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png", "--release"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"describe", "service", serviceName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: update postgres cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/postgrescluster_variation_load_balancer.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: update back to original postgres cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/postgrescluster_original.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: add new service plan to postgres cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/postgrescluster_variation_new_tier.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"remove", "--service-id", cmd.ServiceID})
	err = cmd.RootCmd.Execute()
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
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	serviceName := "mysql cluster" + uuid.NewString()
	// PASS: create mysql cluster service in dev environment
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png", "--release"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	devEnvironmentID := cmd.EnvironmentID
	devProductTierID := cmd.ProductTierID
	require.NotEmpty(devEnvironmentID)
	require.NotEmpty(devProductTierID)

	// PASS: create mysql cluster service in prod environment
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_original.yaml", "--name", serviceName, "--environment", "prod", "--environment-type", "prod", "--release"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	prodEnvironmentID := cmd.EnvironmentID
	prodProductTierID := cmd.ProductTierID
	require.NotEmpty(prodEnvironmentID)
	require.NotEmpty(prodEnvironmentID)
	require.NotEqual(devEnvironmentID, prodEnvironmentID)
	require.NotEqual(devProductTierID, prodProductTierID)

	// PASS: update dev mysql cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_variation_apiparam_image_infra_capability.yaml", "--name", serviceName, "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	require.Equal(devEnvironmentID, cmd.EnvironmentID)
	require.Equal(devProductTierID, cmd.ProductTierID)

	// PASS: update prod mysql cluster service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/variations/mysqlcluster_variation_apiparam_image_infra_capability.yaml", "--name", serviceName, "--environment", "prod", "--environment-type", "prod", "--release-as-preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	require.Equal(prodEnvironmentID, cmd.EnvironmentID)
	require.Equal(prodProductTierID, cmd.ProductTierID)

	cmd.RootCmd.SetArgs([]string{"remove", "--service-id", cmd.ServiceID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}

func Test_build_invalid_file(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"build", "-f", "invalid_file.yaml", "--name", "My Service" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
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
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"build", "--name", "My Service" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
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
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
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
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"remove", "--service-id", cmd.ServiceID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}

func Test_build_create_no_service_logo_url(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--description", "My Service Description"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"remove", "--service-id", cmd.ServiceID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}
