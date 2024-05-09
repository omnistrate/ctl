package utils

import (
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
)

// GetToken returns the authentication token for current user
func GetToken() (string, error) {
	authConfig, err := config.LookupAuthConfig()
	if err != nil {
		return "", err
	}

	return authConfig.Token, nil
}

// GetHost returns the host of the Omnistrate server
func GetHost() string {
	return "api." + utils.GetEnv("ROOT_DOMAIN", "omnistrate.cloud")
}

// GetRootDomain returns the root domain of the Omnistrate server
func GetRootDomain() string {
	return utils.GetEnv("ROOT_DOMAIN", "omnistrate.cloud")
}
