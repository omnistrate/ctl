package dataaccess

import (
	"context"
	"fmt"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	signinapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/signin_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/utils"
	"github.com/pkg/errors"
	goa "goa.design/goa/v3/pkg"
)

func LoginWithPassword(email string, pass string) (token string, err error) {
	signin, err := httpclientwrapper.NewSignin(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return "", err
	}

	request := signinapi.SigninRequest{
		Email:    email,
		Password: commonutils.ToPtr(pass),
	}

	res, err := signin.Signin(context.Background(), &request)
	if err != nil {
		var serviceErr *goa.ServiceError
		ok := errors.As(err, &serviceErr)
		if !ok {
			return
		}

		return "", fmt.Errorf("%s\nDetail: %s", serviceErr.Name, serviceErr.Message)
	}

	token = res.JWTToken
	return
}

func LoginWithIdentityProvider(deviceCode, identityProviderName string) (*signinapi.LoginWithIdentityProviderResult, error) {
	signin, err := httpclientwrapper.NewSignin(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	res, err := signin.LoginWithIdentityProvider(context.Background(), &signinapi.LoginWithIdentityProviderRequest{
		DeviceCode:           commonutils.ToPtr(deviceCode),
		IdentityProviderName: signinapi.IdentityProviderName(identityProviderName),
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}
