package cmd

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
)

func Test_describe_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"describe"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "must provide --service-id")

	rootCmd.SetArgs([]string{"describe", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}

func Test_describe_no_service_logo_url(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--description", "My Service Description"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"describe", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"remove", "--service-id", serviceID})
	err = rootCmd.Execute()
	require.NoError(err)
}
