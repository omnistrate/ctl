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
	secretUpdateExample = `# Update a secret for dev environment
omctl environment secret update dev my-secret new-value

# Update a secret for prod environment
omctl environment secret update prod db-password new-secret123`
)

var secretUpdateCmd = &cobra.Command{
	Use:          "update [environment-type] [secret-name] [secret-value] [flags]",
	Short:        "Update an environment secret",
	Long:         `This command helps you update an existing secret for a specific environment type.`,
	Example:      secretUpdateExample,
	RunE:         runSecretUpdate,
	SilenceUsage: true,
}

func init() {
	secretUpdateCmd.Args = cobra.ExactArgs(3)
}

func runSecretUpdate(cmd *cobra.Command, args []string) error {
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

	// Update the secret (same API as create)
	err = dataaccess.CreateOrUpdateSecret(cmd.Context(), token, environmentType, secretName, secretValue)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	fmt.Printf("Successfully updated secret '%s' for environment type '%s'\n", secretName, environmentType)
	return nil
}