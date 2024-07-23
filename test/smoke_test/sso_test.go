package smoke

import (
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/omnistrate/commons/pkg/utils"
)

func Test_sso(t *testing.T) {
	utils.SmokeTest(t)
	require := require.New(t)

	request := signinapi.LoginWithIdentityProviderRequest{
		AuthorizationCode:    "e2b2c1626a7ae532ce90",
		IdentityProviderName: signinapi.IdentityProviderName("GitHub"),
		RedirectURI:          utils.ToPtr("https://omnistrate.dev/idp-auth"),
	}

	_, err := dataaccess.LoginWithIdentityProvider(request)
	require.Error(err)
	require.Contains(err.Error(), "cannot login with personal email address. Please use your company email address")
}
