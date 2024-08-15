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

	//// PASS: instance list
	//cmd.RootCmd.SetArgs([]string{"instance", "list"})
	//err = cmd.RootCmd.Execute()
	//require.NoError(t, err)

	// PASS: instance list with filters
	cmd.RootCmd.SetArgs([]string{"instance", "list", "--filters", "environment:DEV"})
	err = cmd.RootCmd.Execute()
	require.NoError(t, err)
}
