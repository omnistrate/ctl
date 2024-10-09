package helm

import (
	"github.com/chelnak/ysmrr"
	helmpackageapi "github.com/omnistrate/api-design/v1/pkg/fleet/gen/helm_package_api"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/spf13/cobra"
)

const (
	describeExample = `# Describe the Redis Operator Helm Chart
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
	token, err := config.GetToken()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Describing Helm Chart..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Retrieve Helm Chart
	var helmPackage *helmpackageapi.HelmPackage
	helmPackage, err = dataaccess.DescribeHelmChart(token, chart, version)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully retrieved Helm Chart")

	// Print output
	err = utils.PrintTextTableJsonOutput(output, helmPackage)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	return nil
}
