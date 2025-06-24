package environment

import (
	"fmt"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	secretCreateExample = `# Create a secret for dev environment
omctl environment secret create dev my-secret my-value

# Create a secret for prod environment
omctl environment secret create prod db-password secret123`
)

var secretCreateCmd = &cobra.Command{
	Use:          "create [environment-type] [secret-name] [secret-value] [flags]",
	Short:        "Create or update an environment secret",
	Long:         `This command helps you create or update a secret for a specific environment type.`,
	Example:      secretCreateExample,
	RunE:         runSecretCreate,
	SilenceUsage: true,
}

func init() {
	secretCreateCmd.Args = cobra.ExactArgs(3)
}

func runSecretCreate(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	environmentType := args[0]
	secretName := args[1]
	secretValue := args[2]

	// Validate environment type
	if err := validateEnvironmentType(environmentType); err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Create or update the secret
	err = dataaccess.CreateOrUpdateSecret(cmd.Context(), token, environmentType, secretName, secretValue)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	fmt.Printf("Successfully created/updated secret '%s' for environment type '%s'\n", secretName, environmentType)
	return nil
}