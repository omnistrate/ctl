package environment

import (
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/model"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	secretListExample = `# List secrets for dev environment
omctl environment secret list dev

# List secrets for prod environment with JSON output
omctl environment secret list prod --output json`
)

var secretListCmd = &cobra.Command{
	Use:          "list [environment-type] [flags]",
	Short:        "List environment secrets",
	Long:         `This command helps you list all secrets for a specific environment type.`,
	Example:      secretListExample,
	RunE:         runSecretList,
	SilenceUsage: true,
}

func init() {
	secretListCmd.Args = cobra.ExactArgs(1)
}

func runSecretList(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	environmentType := args[0]
	output, _ := cmd.Flags().GetString("output")

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

	// List secrets
	result, err := dataaccess.ListSecrets(cmd.Context(), token, environmentType)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Format response
	secrets := make([]model.Secret, 0, len(result.GetSecrets()))
	for _, secret := range result.GetSecrets() {
		secrets = append(secrets, model.Secret{
			EnvironmentType: secret.GetEnvironmentType(),
			Name:            secret.GetName(),
			Value:           "[HIDDEN]", // Don't show actual values in list
		})
	}

	secretList := model.SecretList{Secrets: secrets}

	err = utils.PrintTextTableJsonOutput(output, secretList)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}