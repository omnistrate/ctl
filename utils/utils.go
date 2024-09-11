package utils

import (
	_ "embed"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"strings"
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

//go:embed public_key.pem
var publicKey []byte

// GetDefaultServiceAuthPublicKey returns the default public key for environment creation
func GetDefaultServiceAuthPublicKey() string {
	return string(publicKey)
}

func IsProd() bool {
	return GetRootDomain() == "omnistrate.cloud"
}

func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if !strings.ContainsAny(s[:maxLen], " .,:;-") {
		return s[:maxLen] + "..."
	}
	return s[:strings.LastIndexAny(s[:maxLen], " .,:;-!")] + "..."
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

func GetStrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
