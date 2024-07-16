package smoke

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	createaccount "github.com/omnistrate/ctl/cmd/create/account"
	deleteaccount "github.com/omnistrate/ctl/cmd/deletec/account"
	describeaccount "github.com/omnistrate/ctl/cmd/describe/account"
	getaccount "github.com/omnistrate/ctl/cmd/get/account"
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
	createaccount.AccountCmd.SetArgs([]string{"account", awsAccountName, "--aws-account-id", "123456789012", "--aws-bootstrap-role-arn", "arn:aws:iam::123456789012:role/role-name"})
	err = createaccount.AccountCmd.Execute()
	require.NoError(err)

	// PASS: create gcp account
	createaccount.AccountCmd.SetArgs([]string{"account", gcpAccountName, "--gcp-project-id", "project-id", "--gcp-project-number", "project-number", "--gcp-service-account-email", "service-account-email"})
	err = createaccount.AccountCmd.Execute()
	require.NoError(err)

	// PASS: get accounts
	getaccount.AccountCmd.SetArgs([]string{"account"})
	err = getaccount.AccountCmd.Execute()
	require.NoError(err)

	// PASS: get accounts by name
	getaccount.AccountCmd.SetArgs([]string{"account", awsAccountName, gcpAccountName})
	err = getaccount.AccountCmd.Execute()
	require.NoError(err)

	// PASS: describe accounts
	describeaccount.AccountCmd.SetArgs([]string{"account", awsAccountName, gcpAccountName})
	err = describeaccount.AccountCmd.Execute()
	require.NoError(err)

	// PASS: delete accounts
	deleteaccount.AccountCmd.SetArgs([]string{"account", awsAccountName, gcpAccountName})
	err = deleteaccount.AccountCmd.Execute()
	require.NoError(err)
}
