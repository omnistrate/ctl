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
	listExample = `  # List all Helm packages that are saved
  omctl helm list`
)

var listCmd = &cobra.Command{
	Use:          "list [flags]",
	Short:        "List all Helm packages that are saved",
	Long:         `This command helps you list all the Helm packages that are saved.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func init() {
	listCmd.Flags().StringP("output", "o", "text", "Output format (text|table|json)")
}

func runList(cmd *cobra.Command, args []string) error {
	// Get flags
	output, _ := cmd.Flags().GetString("output")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var helmPackageResult *helmpackageapi.ListHelmPackagesResult
	helmPackageResult, err = dataaccess.ListHelmCharts(token)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var jsonData []string
	for _, helmPackage := range helmPackageResult.HelmPackages {
		data, err := json.MarshalIndent(helmPackage, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}

		jsonData = append(jsonData, string(data))
	}

	if len(jsonData) == 0 {
		utils.PrintInfo("No Helm packages found.")
		return nil
	}

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
		fmt.Printf("%+v\n", jsonData)

	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}
