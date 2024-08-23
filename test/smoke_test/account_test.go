package smoke

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/test/testutils"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_account_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	awsAccountName := "aws" + uuid.NewString()
	gcpAccountName := "gcp" + uuid.NewString()

	// FAIL: create aws account
	cmd.RootCmd.SetArgs([]string{"account", "create", awsAccountName, "--aws-account-id", "123456789012"})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "unauthorized: only root users can onboard accounts")

	// FAIL: create gcp account
	cmd.RootCmd.SetArgs([]string{"account", "create", gcpAccountName, "--gcp-project-id", "project-id", "--gcp-project-number", "project-number"})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "unauthorized: only root users can onboard accounts")

	// PASS: list accounts
	cmd.RootCmd.SetArgs([]string{"account", "list"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: list accounts by name
	cmd.RootCmd.SetArgs([]string{"account", "list", "--filter", fmt.Sprintf("name:%s", awsAccountName), "--filter", fmt.Sprintf("name:%s", gcpAccountName)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe accounts
	cmd.RootCmd.SetArgs([]string{"account", "describe", awsAccountName, gcpAccountName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// FAIL: delete accounts
	cmd.RootCmd.SetArgs([]string{"account", "delete", awsAccountName, gcpAccountName})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "account(s) not found")
}
