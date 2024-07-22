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

	// PASS: create aws account
	cmd.RootCmd.SetArgs([]string{"create", "account", awsAccountName, "--aws-account-id", "123456789012"})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "unauthorized: only root users can onboard accounts")

	// PASS: create gcp account
	cmd.RootCmd.SetArgs([]string{"create", "account", gcpAccountName, "--gcp-project-id", "project-id", "--gcp-project-number", "project-number"})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "unauthorized: only root users can onboard accounts")

	// PASS: get accounts
	cmd.RootCmd.SetArgs([]string{"get", "account"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: get accounts by name
	cmd.RootCmd.SetArgs([]string{"get", "account", awsAccountName, gcpAccountName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: describe accounts
	cmd.RootCmd.SetArgs([]string{"describe", "account", awsAccountName, gcpAccountName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: delete accounts
	cmd.RootCmd.SetArgs([]string{"delete", "account", awsAccountName, gcpAccountName})
	err = cmd.RootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "account(s) not found")
}
