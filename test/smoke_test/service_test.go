package smoke_test

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/build"
	deleteservice "github.com/omnistrate/ctl/cmd/deletec/service"
	describeservice "github.com/omnistrate/ctl/cmd/describe/service"
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

	// Get services
	getservice.ServiceCmd.SetArgs([]string{"service"})
	err = getservice.ServiceCmd.Execute()
	require.NoError(err)

	// Get services by name
	getservice.ServiceCmd.SetArgs([]string{"service", serviceName1, serviceName2})
	err = getservice.ServiceCmd.Execute()
	require.NoError(err)

	// Get services by ID
	getservice.ServiceCmd.SetArgs([]string{"service", serviceID1, serviceID2, "--id"})
	err = getservice.ServiceCmd.Execute()
	require.NoError(err)

	// Describe services by name
	describeservice.ServiceCmd.SetArgs([]string{"service", serviceName1, serviceName2})
	err = describeservice.ServiceCmd.Execute()
	require.NoError(err)

	// Describe services by ID
	describeservice.ServiceCmd.SetArgs([]string{"service", serviceID1, serviceID2, "--id"})
	err = describeservice.ServiceCmd.Execute()
	require.NoError(err)

	// Delete services by name
	deleteservice.ServiceCmd.SetArgs([]string{"service", serviceName1, serviceName2})
	err = deleteservice.ServiceCmd.Execute()
	require.NoError(err)

	// Delete services by ID
	deleteservice.ServiceCmd.SetArgs([]string{"service", serviceID1, serviceID2, "--id"})
	err = deleteservice.ServiceCmd.Execute()
	require.Error(err) // Should fail because services were already deleted
}
