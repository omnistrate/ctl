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
	listExample = `# List all Helm packages that are saved
omnistrate helm list`
)

var listCmd = &cobra.Command{
	Use:          "list",
	Short:        "List all Helm packages that are saved.",
	Long:         `This command helps you list all the Helm packages that are saved.`,
	Example:      listExample,
	RunE:         runList,
	SilenceUsage: true,
}

func runList(cmd *cobra.Command, args []string) error {
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

	for _, helmPackage := range helmPackageResult.HelmPackages {
		data, err := json.MarshalIndent(helmPackage, "", "    ")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		fmt.Println(string(data))
	}

	return nil
}
