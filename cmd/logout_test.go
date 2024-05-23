package cmd

import (
	"fmt"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
)

func Test_logout(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	// FAIL: logout without login
	rootCmd.SetArgs([]string{"logout"})
	err = rootCmd.Execute()
	require.Error(err)
	require.Contains(err.Error(), "config file not found")

	// PASS: logout after login
	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"logout"})
	err = rootCmd.Execute()
	require.NoError(err)
}
