package smoke

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd"
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

	// PASS: get accounts
	getaccount.AccountCmd.SetArgs([]string{"account"})
	err = getaccount.AccountCmd.Execute()
	require.NoError(err)

	// PASS: get account with name
	accountName := "BYOA Account"
	getaccount.AccountCmd.SetArgs([]string{"account", accountName})
	err = getaccount.AccountCmd.Execute()
	require.NoError(err)
}
