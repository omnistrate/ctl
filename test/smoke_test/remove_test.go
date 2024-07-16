package smoke

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd"
	"github.com/omnistrate/ctl/cmd/build"
	"github.com/omnistrate/ctl/test/testutils"
	"testing"

	"github.com/google/uuid"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/stretchr/testify/require"
)

func Test_remove_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"build", "-f", "composefiles/postgresql.yaml", "--name", "postgresql" + uuid.NewString(), "--description", "My Service Description", "--service-logo-url", "https://freepnglogos.com/uploads/server-png/server-computer-database-network-vector-graphic-pixabay-31.png"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	cmd.RootCmd.SetArgs([]string{"remove", "--service-id", build.ServiceID})
	err = cmd.RootCmd.Execute()
	require.NoError(err)
}
