package dataaccess

import (
	"context"

	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	"github.com/omnistrate/ctl/config"
	"github.com/omnistrate/ctl/utils"
	openapiclient "github.com/omnistrate/omnistrate-sdk-go/v1"
)

func LoginWithPassword(email string, pass string) (token string, err error) {
	ctx := context.Background()
	request := *openapiclient.NewSigninRequestBody(email)
	request.Password = utils.ToPtr(pass)

	apiClient := getV1Client()
	resp, _, err := apiClient.SigninApiAPI.SigninApiSignin(ctx).SigninRequestBody(request).Execute()
	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	token = resp.JwtToken
	return
}

func LoginWithIdentityProvider(deviceCode, identityProviderName string) (*signinapi.LoginWithIdentityProviderResult, error) {
	signin, err := httpclientwrapper.NewSignin(config.GetHostScheme(), config.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := signin.LoginWithIdentityProvider(context.Background(), &signinapi.LoginWithIdentityProviderRequest{
		DeviceCode:           utils.ToPtr(deviceCode),
		IdentityProviderName: signinapi.IdentityProviderName(identityProviderName),
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
