package common

import (
	"github.com/omnistrate-oss/ctl/cmd/auth/login"
	"github.com/omnistrate-oss/ctl/internal/config"
	"github.com/pkg/errors"
)

func GetTokenWithLogin() (token string, err error) {
	token, err = config.GetToken()
	if err != nil && !errors.Is(err, config.ErrAuthConfigNotFound) && !errors.Is(err, config.ErrConfigFileNotFound) {
		return
	}

	// If token is already present, return it
	if token != "" {
		return
	}

	// Run login command
	err = login.RunLogin(login.LoginCmd, []string{})
	if err != nil {
		return
	}

	token, err = config.GetToken()
	if err != nil {
		return
	}

	return
}
