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

func Test_service_plan_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: create postgresql service
	serviceName := "postgresql" + uuid.NewString()
	cmd.RootCmd.SetArgs([]string{"build", "--file", "composefiles/postgresql.yaml", "--name", serviceName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: release postgresql service plan
	cmd.RootCmd.SetArgs([]string{"service-plan", "release", serviceName, "postgresql"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: set it as preferred
	cmd.RootCmd.SetArgs([]string{"service-plan", "set-default", serviceName, "postgresql", "--version=latest"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: list service plans
	cmd.RootCmd.SetArgs([]string{"service-plan", "list", "--filter", fmt.Sprintf("service_name:%s", serviceName), "--filter", "plan_name:postgresql"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe service plan
	cmd.RootCmd.SetArgs([]string{"service-plan", "describe", serviceName, "postgresql"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: list service plan versions
	cmd.RootCmd.SetArgs([]string{"service-plan", "list-versions", serviceName, "postgresql", "--filter", "service_name:postgresql", "--filter", "plan_name:postgresql"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe service plan version
	cmd.RootCmd.SetArgs([]string{"service-plan", "describe-version", serviceName, "postgresql", "--version=latest"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: delete postgresql service plan
	cmd.RootCmd.SetArgs([]string{"service-plan", "delete", serviceName, "postgresql"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: delete postgresql service
	cmd.RootCmd.SetArgs([]string{"service", "delete", serviceName, "--id=false"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: create postgresql service
	serviceName = "postgresql" + uuid.NewString()
	cmd.RootCmd.SetArgs([]string{"build", "--file", "composefiles/postgresql.yaml", "--name", serviceName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
	serviceId := build.ServiceID
	productTierId := build.ProductTierID

	// PASS: release postgresql service plan
	cmd.RootCmd.SetArgs([]string{"service-plan", "release", "--service-id", serviceId, "--plan-id", productTierId})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: set it as preferred
	cmd.RootCmd.SetArgs([]string{"service-plan", "set-default", "--service-id", serviceId, "--plan-id", productTierId, "--version=latest"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: list service plans
	cmd.RootCmd.SetArgs([]string{"service-plan", "list", "--filter", fmt.Sprintf("service_id:%s", serviceId), "--filter", fmt.Sprintf("plan_id:%s", productTierId)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe service plan
	cmd.RootCmd.SetArgs([]string{"service-plan", "describe", "--service-id", serviceId, "--plan-id", productTierId})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: list service plan versions
	cmd.RootCmd.SetArgs([]string{"service-plan", "list-versions", "--service-id", serviceId, "--plan-id", productTierId})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe service plan version
	cmd.RootCmd.SetArgs([]string{"service-plan", "describe-version", "--service-id", serviceId, "--plan-id", productTierId, "--version=preferred"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: delete postgresql service plan
	cmd.RootCmd.SetArgs([]string{"service-plan", "delete", "--service-id", serviceId, "--plan-id", productTierId})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: delete postgresql service
	cmd.RootCmd.SetArgs([]string{"service", "delete", serviceName, "--id=false"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}
