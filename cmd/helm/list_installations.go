package helm

import (
	"encoding/json"
	"fmt"
	helmpackageapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/helm_package_api"
	commonutils "github.com/omnistrate/commons/pkg/utils"
	"github.com/omnistrate/ctl/dataaccess"
	"github.com/omnistrate/ctl/utils"
	"github.com/spf13/cobra"
)

const (
	listInstallationsExample = `# List all Helm Packages and the Kubernetes clusters that they are installed on
omnistrate helm list-installations --host-cluster-id=[host-cluster-id]`
)

type helmPackageInstallationIntermediate struct {
	ChartName     string
	ChartVersion  string
	Namespace     string
	HostClusterID string
	Status        string
}

var listInstallationsCmd = &cobra.Command{
	Use:          "list-installations --host-cluster-id=[host-cluster-id]",
	Short:        "List all Helm Packages and the Kubernetes clusters that they are installed on",
	Long:         `This command helps you list all the Helm Packages and the Kubernetes clusters that they are installed on.`,
	Example:      listInstallationsExample,
	RunE:         runListInstallations,
	SilenceUsage: true,
}

func init() {
	saveCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	listInstallationsCmd.Flags().String("host-cluster-id", "", "Host cluster ID")
	listInstallationsCmd.Flags().StringP("output", "o", "text", "Output format (text|json)")
}

func runListInstallations(cmd *cobra.Command, args []string) error {
	// Get flags
	hostClusterID, _ := cmd.Flags().GetString("host-cluster-id")
	output, _ := cmd.Flags().GetString("output")

	// Validate user is currently logged in
	token, err := utils.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var hostClusterIDReq *helmpackageapi.HostClusterID
	var helmPackageResult *helmpackageapi.ListHelmPackageInstallationsResult

	if len(hostClusterID) > 0 {
		hostClusterIDReq = commonutils.ToPtr(helmpackageapi.HostClusterID(hostClusterID))
	}

	helmPackageResult, err = dataaccess.ListHelmChartInstallations(token, hostClusterIDReq)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	var jsonData []string
	for _, helmPackageInstallation := range helmPackageResult.HelmPackageInstallations {
		// Convert HelmPackageInstallation to intermediate struct
		intermediate := helmPackageInstallationIntermediate{
			ChartName:     helmPackageInstallation.HelmPackage.ChartName,
			ChartVersion:  helmPackageInstallation.HelmPackage.ChartVersion,
			Namespace:     helmPackageInstallation.HelmPackage.Namespace,
			HostClusterID: string(helmPackageInstallation.HostClusterID),
			Status:        helmPackageInstallation.Status,
		}
		data, err := json.MarshalIndent(intermediate, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		jsonData = append(jsonData, string(data))
	}

	if len(jsonData) == 0 {
		utils.PrintInfo("No Helm package installations found.")
		return nil
	}

	switch output {
	case "text":
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
