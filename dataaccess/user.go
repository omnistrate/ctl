package dataaccess

import (
	"context"
	"github.com/omnistrate/api-design/pkg/httpclientwrapper"
	usersapi "github.com/omnistrate/api-design/v1/pkg/registration/gen/users_api"
	"github.com/omnistrate/ctl/utils"
)

func DescribeUser(token string) (*usersapi.DescribeUserResult, error) {
	user, err := httpclientwrapper.NewUser(utils.GetHostScheme(), utils.GetHost())
	if err != nil {
		return nil, err
	}

	request := usersapi.DescribeUserRequest{
		Token: token,
	}

	res, err := user.DescribeUser(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
