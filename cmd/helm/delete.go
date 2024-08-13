package helm

import (
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	deleteExample = `# Delete a Helm package
omnistrate helm delete redis --version=20.0.1`
)

var deleteCmd = &cobra.Command{
	Use:          "delete chart --version=[version]",
	Short:        "Delete a Helm package for your service.",
	Long:         `This command helps you delete the templates for your helm packages.`,
	Example:      deleteExample,
	RunE:         runDelete,
	SilenceUsage: true,
}

func init() {
	deleteCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	deleteCmd.Flags().String("version", "", "Helm Chart version")

	err := deleteCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

func runDelete(cmd *cobra.Command, args []string) error {
	chart := args[0]
	version, _ := cmd.Flags().GetString("version")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	err = dataaccess.DeleteHelmChart(token, chart, version)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	utils.PrintSuccess("Helm package deleted successfully")

	return nil
}
