package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	"github.com/omnistrate/ctl/utils"
)

func LoginWithPassword(request signinapi.SigninRequest) (*signinapi.SigninResult, error) {
	signin, err := httpclientwrapper.NewSignin(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := signin.Signin(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func LoginWithIdentityProvider(request signinapi.LoginWithIdentityProviderRequest) (*signinapi.LoginWithIdentityProviderResult, error) {
	signin, err := httpclientwrapper.NewSignin(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := signin.LoginWithIdentityProvider(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
