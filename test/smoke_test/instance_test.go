package smoke

import (
	"fmt"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/cmd"
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

	// PASS: instance create
	cmd.RootCmd.SetArgs([]string{"instance", "create",
		"--service=mySQL",
		"--environment=dev",
		"--plan=mysql",
		"--version=latest",
		"--resource=mySQL",
		"--cloud-provider=aws",
		"--region=ca-central-1",
		"--param", `{"databaseName":"default","password":"a_secure_password","rootPassword":"a_secure_root_password","username":"user"},"productTierVersion":"1.0"}`})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// PASS: instance list
	cmd.RootCmd.SetArgs([]string{"instance", "list"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)

	// PASS: instance list with filters
	cmd.RootCmd.SetArgs([]string{"instance", "list", "-f", "environment:DEV,cloud_provider:gcp", "-f", "environment:DEV,cloud_provider:aws"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)
}
