package config

import (
	_ "embed"

	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// GetToken returns the authentication token for current user
func GetToken() (string, error) {
	authConfig, err := LookupAuthConfig()
	if err != nil {
		return "", err
	}

	return authConfig.Token, nil
}

// GetHost returns the host of the Omnistrate server
func GetHost() string {
	return utils.GetEnv("OMNISTRATE_HOST", "api"+"."+GetRootDomain())
}

// GetRootDomain returns the root domain of the Omnistrate server
func GetRootDomain() string {
	return utils.GetEnv("OMNISTRATE_ROOT_DOMAIN", "omnistrate.cloud")
}

// GetHostScheme returns the scheme of the Omnistrate server
func GetHostScheme() string {
	return utils.GetEnv("OMNISTRATE_HOST_SCHEME", "https")
}

func GetDebug() bool {
	return utils.GetEnvAsBoolean("OMNISTRATE_DEBUG", "false")
}

//go:embed public_key.pem
var publicKey []byte

// GetDefaultServiceAuthPublicKey returns the default public key for environment creation
func GetDefaultServiceAuthPublicKey() string {
	return string(publicKey)
}

func IsProd() bool {
	return GetRootDomain() == "omnistrate.cloud"
}

func CleanupArgsAndFlags(cmd *cobra.Command, args *[]string) {
	// Clean up flags
	cmd.Flags().VisitAll(
		func(f *pflag.Flag) {
			_ = cmd.Flags().Set(f.Name, f.DefValue)
		})

	// Clean up arguments by resetting the slice to nil or an empty slice
	*args = nil
}
