package dataaccess

import (
	"context"

	openapiclient "github.com/omnistrate-oss/omnistrate-sdk-go/v1"
	"github.com/omnistrate/ctl/internal/utils"
)

func LoginWithPassword(ctx context.Context, email string, pass string) (string, error) {
	request := *openapiclient.NewSigninRequestBody(email)
	request.Password = utils.ToPtr(pass)

	apiClient := getV1Client()
	resp, r, err := apiClient.SigninApiAPI.SigninApiSignin(ctx).SigninRequestBody(request).Execute()

	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	r.Body.Close()
	return resp.JwtToken, nil
}

func LoginWithIdentityProvider(ctx context.Context, deviceCode, identityProviderName string) (string, error) {
	request := *openapiclient.NewLoginWithIdentityProviderRequestBody(identityProviderName)
	request.DeviceCode = utils.ToPtr(deviceCode)

	apiClient := getV1Client()
	resp, r, err := apiClient.SigninApiAPI.SigninApiLoginWithIdentityProvider(ctx).LoginWithIdentityProviderRequestBody(request).Execute()

	err = handleV1Error(err)
	if err != nil {
		return "", err
	}

	r.Body.Close()
	return resp.JwtToken, nil
}
