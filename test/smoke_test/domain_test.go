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

func Test_domain_basic(t *testing.T) {
	utils.SmokeTest(t)

	require := require.New(t)
	defer testutils.Cleanup()

	var err error

	testEmail, testPassword, err := testutils.GetSmokeTestAccount()
	require.NoError(err)
	cmd.RootCmd.SetArgs([]string{"login", fmt.Sprintf("--email=%s", testEmail), fmt.Sprintf("--password=%s", testPassword)})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	devDomainName := "dev" + uuid.NewString()
	devDomain := "domain" + uuid.NewString() + ".dev"
	prodDomainName := "prod" + uuid.NewString()
	prodDomain := "domain" + uuid.NewString() + ".prod"

	// create dev domain
	cmd.RootCmd.SetArgs([]string{"create", "domain", devDomainName, "--env", "dev", "--domain", devDomain})
	err = cmd.RootCmd.Execute()
	if err != nil {
		require.Contains(err.Error(), "saas portal does not exist for environment type")
	}

	// create prod domain
	cmd.RootCmd.SetArgs([]string{"create", "domain", prodDomainName, "--env", "prod", "--domain", prodDomain})
	err = cmd.RootCmd.Execute()
	if err != nil {
		require.Contains(err.Error(), "saas portal does not exist for environment type")
	}

	// PASS: get domains
	cmd.RootCmd.SetArgs([]string{"get", "domain"})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// PASS: get domains by name
	cmd.RootCmd.SetArgs([]string{"get", "domain", devDomainName, prodDomainName})
	err = cmd.RootCmd.Execute()
	require.NoError(err)

	// delete domains
	cmd.RootCmd.SetArgs([]string{"delete", "domain", devDomainName, prodDomainName})
	err = cmd.RootCmd.Execute()
	if err != nil {
		require.Contains(err.Error(), "domain(s) not found")
	}
}
