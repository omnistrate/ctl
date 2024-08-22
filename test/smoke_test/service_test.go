package smoke

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/build"
	"github.com/omnistrate/ctl/test/testutils"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_service_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	serviceName1 := "postgresql" + uuid.NewString()
	serviceName2 := "postgresql" + uuid.NewString()

	// Build service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--name", serviceName1, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	serviceID1 := build.ServiceID

	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--name", serviceName2, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	serviceID2 := build.ServiceID

	// List services
	cmd.RootCmd.SetArgs([]string{"service", "list"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// List services by name
	cmd.RootCmd.SetArgs([]string{"service", "list", "--filter", fmt.Sprintf("name:%s", serviceName1), "--filter", fmt.Sprintf("name:%s", serviceName2)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Get services by ID
	cmd.RootCmd.SetArgs([]string{"service", "list", "--filter", fmt.Sprintf("id:%s", serviceID1), "--filter", fmt.Sprintf("id:%s", serviceID2), "--id"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Describe services by name
	cmd.RootCmd.SetArgs([]string{"service", "describe", serviceName1, serviceName2})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Describe services by ID
	cmd.RootCmd.SetArgs([]string{"service", "describe", serviceID1, serviceID2, "--id"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Delete service by name
	cmd.RootCmd.SetArgs([]string{"service", "delete", serviceName1})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Delete service by ID
	cmd.RootCmd.SetArgs([]string{"service", "delete", serviceID2, "--id"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}
