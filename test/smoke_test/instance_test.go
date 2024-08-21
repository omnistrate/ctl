package smoke

import (
	"fmt"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/instance"
	"github.com/omnistrate/ctl/test/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInstanceBasic(t *testing.T) {
	utils.SmokeTest(t)

	defer testutils.Cleanup()

	// PASS: login
	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(t, err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// PASS: create instance 1 with param
	cmd.RootCmd.SetArgs([]string{"instance", "create",
		"--service=mysql",
		"--environment=dev",
		"--plan=mysql",
		"--version=latest",
		"--resource=mySQL",
		"--cloud-provider=aws",
		"--region=ca-central-1",
		"--param", `{"databaseName":"default","password":"a_secure_password","rootPassword":"a_secure_root_password","username":"user"}`})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)
	instanceID1 := instance.InstanceID
	require.NotEmpty(t, instanceID1)

	// PASS: create instance 2 with param file
	cmd.RootCmd.SetArgs([]string{"instance", "create",
		"--service=mysql",
		"--environment=dev",
		"--plan=mysql",
		"--version=latest",
		"--resource=mySQL",
		"--cloud-provider=aws",
		"--region=ca-central-1",
		"--param-file", "paramfiles/instance_create_param.json"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)
	instanceID2 := instance.InstanceID
	require.NotEmpty(t, instanceID2)

	// PASS: describe instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "describe", instanceID1})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// PASS: describe instance 2
	cmd.RootCmd.SetArgs([]string{"instance", "describe", instanceID2})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// TODO: Uncomment the following tests once dev search resource instance is fixed
	//// PASS: instance list
	//cmd.RootCmd.SetArgs([]string{"instance", "list"})
	//err = cmd.RootCmd.Execute()
	//require.NoError(t, err)
	//
	//// PASS: instance list with filters
	//cmd.RootCmd.SetArgs([]string{"instance", "list", "-f", "environment:DEV,cloud_provider:gcp", "-f", "environment:DEV,cloud_provider:aws"})
	//err = cmd.RootCmd.Execute()
	//require.NoError(t, err)

	// PASS: delete instance 1
	cmd.RootCmd.SetArgs([]string{"instance", "delete", instanceID1, "--yes"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// PASS: delete instance 2
	cmd.RootCmd.SetArgs([]string{"instance", "delete", instanceID2, "--yes"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)
}
