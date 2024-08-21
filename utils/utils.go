package utils

import (
	_ "embed"
	"github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/config"
	"github.com/spf13/cobra"
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

func CombineSubCmdExamples(root *cobra.Command) (example string) {
	for _, cmd := range root.Commands() {
		example += cmd.Example + "\n\n"
	}
	return
}

func TruncateString(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if !strings.ContainsAny(s[:max], " .,:;-") {
		return s[:max] + "..."
	}
	return s[:strings.LastIndexAny(s[:max], " .,:;-!")] + "..."
}
