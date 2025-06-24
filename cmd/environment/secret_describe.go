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
	secretDescribeExample = `# Describe a secret in dev environment
omctl environment secret describe dev my-secret

# Describe a secret with JSON output
omctl environment secret describe prod db-password --output json`
)

var secretDescribeCmd = &cobra.Command{
	Use:          "describe [environment-type] [secret-name] [flags]",
	Short:        "Describe an environment secret",
	Long:         `This command helps you describe a specific secret for an environment type.`,
	Example:      secretDescribeExample,
	RunE:         runSecretDescribe,
	SilenceUsage: true,
}

func init() {
	secretDescribeCmd.Args = cobra.ExactArgs(2)
}

func runSecretDescribe(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	environmentType := args[0]
	secretName := args[1]
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

	// Get secret
	result, err := dataaccess.GetSecret(cmd.Context(), token, environmentType, secretName)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Format response
	secret := model.Secret{
		EnvironmentType: result.GetEnvironmentType(),
		Name:            result.GetName(),
		Value:           result.GetValue(),
	}

	err = utils.PrintTextTableJsonOutput(output, secret)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}