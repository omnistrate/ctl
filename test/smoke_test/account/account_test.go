package account

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/test/testutils"

	"github.com/stretchr/testify/require"
)

func Test_account_basic(t *testing.T) {
	testutils.SmokeTest(t)

	ctx := context.TODO()

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	awsAccountName := "aws" + uuid.NewString()
	gcpAccountName := "gcp" + uuid.NewString()
	awsAccountNumber := "903112026573" // Use acc # that is not onboarded yet, otherwise it will fail with different error

	// FAIL: create aws account
	cmd.RootCmd.SetArgs([]string{"account", "create", awsAccountName, "--aws-account-id", awsAccountNumber})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), "unauthorized: only root users can onboard accounts")

	// FAIL: create gcp account
	cmd.RootCmd.SetArgs([]string{"account", "create", gcpAccountName, "--gcp-project-id", "project-id", "--gcp-project-number", "project-number"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), "unauthorized: only root users can onboard accounts")

	// PASS: list accounts
	cmd.RootCmd.SetArgs([]string{"account", "list"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: list accounts by name
	cmd.RootCmd.SetArgs([]string{"account", "list", "--filter", fmt.Sprintf("name:%s", awsAccountName), "--filter", fmt.Sprintf("name:%s", gcpAccountName)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: describe account
	cmd.RootCmd.SetArgs([]string{"account", "describe", awsAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), "account not found")

	cmd.RootCmd.SetArgs([]string{"account", "describe", gcpAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), "account not found")

	// FAIL: delete account
	cmd.RootCmd.SetArgs([]string{"account", "delete", awsAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), "account not found")

	cmd.RootCmd.SetArgs([]string{"account", "delete", gcpAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.Error(err)
	require.Contains(err.Error(), "account not found")
}
