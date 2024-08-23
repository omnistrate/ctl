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
	describeExample = `  # Describe the Redis Operator Helm Chart
  omctl helm describe redis --version=20.0.1`
)

var describeCmd = &cobra.Command{
	Use:          "describe chart --version=[version]",
	Short:        "Describe a Helm Chart for your service",
	Long:         `This command helps you describe the templates for your helm charts.`,
	Example:      describeExample,
	RunE:         runDescribe,
	SilenceUsage: true,
}

func init() {
	describeCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	describeCmd.Flags().String("version", "", "Helm Chart version")
	describeCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")

	err := describeCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Get flags
	chart := args[0]
	version, _ := cmd.Flags().GetString("version")
	output, _ := cmd.Flags().GetString("output")

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

	var jsonData []string
	data, err := json.MarshalIndent(helmPackage, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	jsonData = append(jsonData, string(data))

	switch output {
	case "text":
		err = utils.PrintText(jsonData)
		if err != nil {
			return err
		}
	case "table":
		err = utils.PrintTable(jsonData)
		if err != nil {
			return err
		}
	case "json":
		fmt.Println(string(data))

	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
