package helm

import (
	"encoding/json"
	"fmt"
	helmpackageapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/helm_package_api"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

const (
	listExample = `# List all Helm packages that are saved
omnistrate helm list`
)

var listCmd = &cobra.Command{
	Use:          "list",
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
		PrintTable(helmPackageResult.HelmPackages)
	case "table":
		var tableWriter *utils.Table
		if tableWriter, err = utils.NewTableFromJSONTemplate(json.RawMessage(jsonData[0])); err != nil {
			// Just print the JSON directly and return
			fmt.Printf("%+v\n", jsonData)
			return err
		}

		for _, data := range jsonData {
			if err = tableWriter.AddRowFromJSON(json.RawMessage(data)); err != nil {
				// Just print the JSON directly and return
				fmt.Printf("%+v\n", jsonData)
				return err
			}
		}

		tableWriter.Print()

	case "json":
		fmt.Printf("%+v\n", jsonData)

	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}

	return nil
}

func PrintTable(res []*helmpackageapi.HelmPackage) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	_, err := fmt.Fprintln(w, "Chart Name\tChart Version\tNamespace\tRepo URL\tValues")
	if err != nil {
		return
	}

	for _, r := range res {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			r.ChartName,
			r.ChartVersion,
			r.Namespace,
			r.RepoURL,
			r.Values,
		)
		if err != nil {
			return
		}
	}

	err = w.Flush()
	if err != nil {
		utils.PrintError(err)
	}
}
