package smoke

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/omnistrate/ctl/cmd"
	createdomain "github.com/omnistrate/ctl/cmd/create/domain"
	deletedomain "github.com/omnistrate/ctl/cmd/deletec/domain"
	getdomain "github.com/omnistrate/ctl/cmd/get/domain"
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
	prodDomainName := "prod" + uuid.NewString()

	// PASS: create dev domain
	createdomain.DomainCmd.SetArgs([]string{devDomainName, "--env", "dev", "--domain", "domain.dev"})
	err = createdomain.DomainCmd.Execute()
	require.NoError(err)

	// PASS: create prod domain
	createdomain.DomainCmd.SetArgs([]string{prodDomainName, "--env", "prod", "--domain", "domain.prod"})
	err = createdomain.DomainCmd.Execute()
	require.NoError(err)

	// PASS: get domains
	getdomain.DomainCmd.SetArgs([]string{})
	err = getdomain.DomainCmd.Execute()
	require.NoError(err)

	// PASS: get domains by name
	getdomain.DomainCmd.SetArgs([]string{devDomainName, prodDomainName})
	err = getdomain.DomainCmd.Execute()
	require.NoError(err)

	// PASS: delete domains
	deletedomain.DomainCmd.SetArgs([]string{devDomainName, prodDomainName})
	err = deletedomain.DomainCmd.Execute()
	require.NoError(err)
}
