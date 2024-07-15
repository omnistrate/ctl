package smoke

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/get/account"
	"github.com/omnistrate/ctl/test/testutils"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_get_account_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	account.AccountCmd.SetArgs([]string{"account"})
	err = account.AccountCmd.Execute()
	require.NoError(err)

	accountName := "BYOA Account"
	account.AccountCmd.SetArgs([]string{"account", accountName})
	err = account.AccountCmd.Execute()
	require.NoError(err)
}
