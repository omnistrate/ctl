package smoke_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	deleteservice "github.com/omnistrate/ctl/cmd/deletec/service"
	getservice "github.com/omnistrate/ctl/cmd/get/service"
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

	serviceName := "postgresql" + uuid.NewString()

	// Build service
	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--name", serviceName, "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// Get service
	getservice.ServiceCmd.SetArgs([]string{"service"})
	err = getservice.ServiceCmd.Execute()
	require.NoError(err)

	// Get service by name
	getservice.ServiceCmd.SetArgs([]string{"service", serviceName})
	err = getservice.ServiceCmd.Execute()
	require.NoError(err)

	// Delete service
	deleteservice.ServiceCmd.SetArgs([]string{"service", serviceName})
	err = deleteservice.ServiceCmd.Execute()
	require.NoError(err)
}
