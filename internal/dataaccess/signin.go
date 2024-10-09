package dataaccess

import (
	"context"

	"github.com/omnistrate/ctl/internal/utils"
	openapiclient "github.com/omnistrate/omnistrate-sdk-go/v1"
)

func LoginWithPassword(ctx context.Context, email string, pass string) (token string, err error) {
	request := *openapiclient.NewSigninRequestBody(email)
	request.Password = utils.ToPtr(pass)

	apiClient := getV1Client()
	resp, _, err := apiClient.SigninApiAPI.SigninApiSignin(ctx).SigninRequestBody(request).Execute()
	err = handleV1Error(err)

	return resp.JwtToken, nil
}

func LoginWithIdentityProvider(ctx context.Context, deviceCode, identityProviderName string) (token string, err error) {
	request := *openapiclient.NewLoginWithIdentityProviderRequestBody(identityProviderName)
	request.DeviceCode = utils.ToPtr(deviceCode)

	apiClient := getV1Client()
	resp, _, err := apiClient.SigninApiAPI.SigninApiLoginWithIdentityProvider(ctx).LoginWithIdentityProviderRequestBody(request).Execute()
	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	return resp.JwtToken, nil
}
