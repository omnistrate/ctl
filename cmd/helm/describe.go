package helm

import (
	"encoding/json"
	"fmt"
	helmpackageapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/helm_package_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe the Redis Operator Helm Chart
omnistrate helm describe redis --version=20.0.1`
)

var describeCmd = &cobra.Command{
	Use:          "describe chart --version=[version]",
	Short:        "Describe a Helm Chart for your service.",
	Long:         `This command helps you describe the templates for your helm charts.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	describeCmd.Flags().String("version", "", "Helm Chart version")

	err := describeCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Get flags
	chart := args[0]
	version, _ := cmd.Flags().GetString("version")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var helmPackage *helmpackageapi.HelmPackage
	helmPackage, err = dataaccess.DescribeHelmChart(token, chart, version)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	data, err := json.MarshalIndent(helmPackage, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	fmt.Println(string(data))

	return nil
}
