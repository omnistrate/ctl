package cmd

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/testutils"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_list_basic(t *testing.T) {
	require := require.New(t)
	defer testutils.Cleanup()

	err := os.Setenv("ROOT_DOMAIN", "omnistrate.dev")
	require.NoError(err)

	testEmail, testPassword := testutils.GetTestAccount()
	rootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"build", "-f", "../composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = rootCmd.Execute()
	require.NoError(err)

	rootCmd.SetArgs([]string{"list"})
	err = rootCmd.Execute()
	require.NoError(err)
}
