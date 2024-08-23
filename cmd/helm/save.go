package helm

import (
	"encoding/json"
	"fmt"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
	"os"
)

const (
	saveExample = `  # Install the Redis Operator Helm Chart
  omctl helm save redis --repo-url=https://charts.bitnami.com/bitnami --version=20.0.1 --namespace=redis-operator`
)

var saveCmd = &cobra.Command{
	Use:          "save chart --repo-url=[repo-url] --version=[version] --namespace=[namespace] --values-file=[values-file]",
	Short:        "Save a Helm Chart for your service",
	Long:         `This command helps you save the templates for your helm charts.`,
	Example:      saveExample,
	RunE:         runSave,
	SilenceUsage: true,
}

func init() {
	saveCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	saveCmd.Flags().String("repo-url", "", "Helm Chart repository URL")
	saveCmd.Flags().String("version", "", "Helm Chart version")
	saveCmd.Flags().String("namespace", "", "Helm Chart namespace")
	saveCmd.Flags().String("values-file", "", "Helm Chart values file containing custom values defined as a JSON")
	saveCmd.Flags().StringP("output", "o", "text", "Output format (text|json)")

	err := saveCmd.MarkFlagRequired("repo-url")
	if err != nil {
		return
	}

	err = saveCmd.MarkFlagRequired("version")
	if err != nil {
		return
	}

	err = saveCmd.MarkFlagRequired("namespace")
	if err != nil {
		return
	}
}

func runSave(cmd *cobra.Command, args []string) error {
	// Get flags
	chart := args[0]
	repoURL, _ := cmd.Flags().GetString("repo-url")
	version, _ := cmd.Flags().GetString("version")
	namespace, _ := cmd.Flags().GetString("namespace")
	valuesFile, _ := cmd.Flags().GetString("values-file")
	output, _ := cmd.Flags().GetString("output")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var values map[string]any
	if len(valuesFile) > 0 {
		// Read Values file as a JSON
		if _, err = os.Stat(valuesFile); os.IsNotExist(err) {
			err = fmt.Errorf("can't load values file from non existent path: %s", valuesFile)
			utils.PrintError(err)
			return err
		}

		var data []byte
		if data, err = os.ReadFile(valuesFile); err != nil {
			utils.PrintError(err)
			return err
		}

		if err = json.Unmarshal(data, &values); err != nil {
			utils.PrintError(err)
			return err
		}
	}

	// Save Helm Chart
	helmPackage, err := dataaccess.SaveHelmChart(token, chart, version, namespace, repoURL, values)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var data []byte
	data, err = json.MarshalIndent(helmPackage, "", "    ")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	switch output {
	case "text":
		utils.PrintSuccess("Helm Chart saved successfully")

		var tableWriter *utils.Table
		if tableWriter, err = utils.NewTableFromJSONTemplate(data); err != nil {
			// Just print the JSON directly and return
			fmt.Println(string(data))
			return err
		}

		if err = tableWriter.AddRowFromJSON(data); err != nil {
			// Just print the JSON directly and return
			fmt.Println(string(data))
			return err
		}

		tableWriter.Print()

	case "json":
		fmt.Println(string(data))

	default:
		err = fmt.Errorf("unsupported output format: %s", output)
		utils.PrintError(err)
		return err
	}
	return nil
}
