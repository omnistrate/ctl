package serviceorchestration

import (
	"encoding/base64"
	"os"

	"github.com/chelnak/ysmrr"
	"github.com/omnistrate/ctl/cmd/common"
	"github.com/omnistrate/ctl/internal/config"
	"github.com/omnistrate/ctl/internal/dataaccess"
	"github.com/omnistrate/ctl/internal/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	createExample = `# Create a service orchestration deployment from a DSL file
omctl service-orchestration create --dsl-file /path/to/dsl.yaml`
)

var ServiceOrchestrationID string

var createCmd = &cobra.Command{
	Use:          "create --dsl-file=[file-path]",
	Short:        "Create a service orchestration deployment",
	Long:         `This command helps you create a service orchestration deployment, coordinating the creation of multiple services.`,
	Example:      createExample,
	RunE:         runCreate,
	SilenceUsage: true,
}

func init() {
	createCmd.Flags().String("dsl-file", "", "Yaml file containing DSL for service orchestration deployment")

	if err := createCmd.MarkFlagRequired("dsl-file"); err != nil {
		return
	}

	createCmd.Args = cobra.NoArgs // Require no arguments
}

func runCreate(cmd *cobra.Command, args []string) error {
	defer config.CleanupArgsAndFlags(cmd, &args)

	// Retrieve flags
	dslFilePath, err := cmd.Flags().GetString("dsl-file")
	if err != nil {
		utils.PrintError(err)
		return err
	}
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Validate user login
	token, err := common.GetTokenWithLogin()
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Initialize spinner if output is not JSON
	var sm ysmrr.SpinnerManager
	var spinner *ysmrr.Spinner
	if output != "json" {
		sm = ysmrr.NewSpinnerManager()
		msg := "Creating service orchestration..."
		spinner = sm.AddSpinner(msg)
		sm.Start()
	}

	// Read DSL file
	dslFileContent, err := readDslFile(dslFilePath)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	orchestration, err := dataaccess.CreateServicesOrchestration(
		cmd.Context(),
		token,
		dslFileContent)
	if err != nil {
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	if orchestration.Id == nil {
		err = errors.New("failed to create service orchestration")
		utils.HandleSpinnerError(spinner, sm, err)
		return err
	}

	utils.HandleSpinnerSuccess(spinner, sm, "Successfully service orchestration")

	// Search for the orchestration
	searchRes, err := dataaccess.DescribeServicesOrchestration(
		cmd.Context(),
		token,
		*orchestration.Id,
	)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	// Print output
	if err = utils.PrintTextTableJsonOutput(output, searchRes); err != nil {
		return err
	}

	return nil
}

// Helper functions

func readDslFile(filePath string) (base64FileContent string, err error) {
	// Read parameters from file if provided
	if filePath == "" {
		err = errors.New("dsl file path is empty")
		return
	}
	var fileContent []byte
	fileContent, err = os.ReadFile(filePath)
	if err != nil {
		return
	}
	// return base64 encoded file content
	base64FileContent = base64.StdEncoding.EncodeToString(fileContent)
	return
}
