package account

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

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
	azureAccountName := "azure" + uuid.NewString()

	rand.Seed(time.Now().UnixNano())
	num := rand.Int63n(9000000000) + 1000000000 // 10 digit number

	// PASS: create aws account
	cmd.RootCmd.SetArgs([]string{"account", "create", awsAccountName, "--aws-account-id", fmt.Sprintf("%d", num)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: create gcp account
	cmd.RootCmd.SetArgs([]string{"account", "create", gcpAccountName, "--gcp-project-id", fmt.Sprintf("project-id-%d", num), "--gcp-project-number", "project-number"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: create azure account
	cmd.RootCmd.SetArgs([]string{"account", "create", azureAccountName, "--azure-subscription-id", fmt.Sprintf("12345678-1234-1234-1234-%d", num), "--azure-tenant-id", "87654321-4321-4321-4321-210987654321"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: list accounts
	cmd.RootCmd.SetArgs([]string{"account", "list"})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: list accounts by name
	cmd.RootCmd.SetArgs([]string{"account", "list", "--filter", fmt.Sprintf("name:%s", awsAccountName), "--filter", fmt.Sprintf("name:%s", gcpAccountName), "--filter", fmt.Sprintf("name:%s", azureAccountName)})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: describe account
	cmd.RootCmd.SetArgs([]string{"account", "describe", awsAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"account", "describe", gcpAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"account", "describe", azureAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	// PASS: delete account
	cmd.RootCmd.SetArgs([]string{"account", "delete", awsAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"account", "delete", gcpAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"account", "delete", azureAccountName})
	err = cmd.RootCmd.ExecuteContext(ctx)
	require.NoError(err)
}
