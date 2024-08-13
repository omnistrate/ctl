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

var listInstallationsCmd = &cobra.Command{
	Use:          "list-installations --host-cluster-id=[host-cluster-id]",
	Short:        "List all Helm Packages and the Kubernetes clusters that they are installed on.",
	Long:         `This command helps you list all the Helm Packages and the Kubernetes clusters that they are installed on.`,
	Example:      listInstallationsExample,
	RunE:         runListInstallations,
	SilenceUsage: true,
}

func init() {
	saveCmd.Args = cobra.ExactArgs(1) // Require exactly one argument

	listInstallationsCmd.Flags().String("host-cluster-id", "", "Host cluster ID")
}

func runListInstallations(cmd *cobra.Command, args []string) error {
	// Get flags
	hostClusterID, _ := cmd.Flags().GetString("host-cluster-id")

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

	for _, helmPackageInstallation := range helmPackageResult.HelmPackageInstallations {
		data, err := json.MarshalIndent(helmPackageInstallation, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
