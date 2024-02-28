package utils

import "github.com/omnistrate/ctl/config"

func GetToken() (string, error) {
	authConfig, err := config.LookupAuthConfig()
	if err != nil {
		return "", err
	}

	return authConfig.Token, nil
}
